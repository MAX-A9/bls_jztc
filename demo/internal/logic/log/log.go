package log

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	v1 "demo/api/log/v1"
	"demo/internal/consts"
	"demo/internal/middleware"
	"demo/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// 常量定义
const (
	// 日志保存目录
	logDir = "resource/logs/operation"
	// 单个日志文件最大行数
	maxLogLines = 1000
	// 最大保留日志文件数
	maxLogFiles = 30
)

// 互斥锁，用于保护日志文件的并发写入
var logMutex sync.Mutex

// logImpl 日志服务实现
type logImpl struct{}

// New 创建一个日志服务实例
func New() service.LogService {
	// 确保日志目录存在
	if !gfile.Exists(logDir) {
		if err := gfile.Mkdir(logDir); err != nil {
			g.Log().Fatal(context.Background(), "创建日志目录失败:", err)
		}
	}
	return &logImpl{}
}

// List 获取操作日志列表
func (s *logImpl) List(ctx context.Context, req *v1.OperationLogListReq) (res *v1.OperationLogListRes, err error) {
	res = &v1.OperationLogListRes{
		List:  make([]v1.OperationLogListItem, 0),
		Total: 0,
		Page:  req.Page,
	}

	// 获取所有日志文件
	logFiles, err := s.getLogFiles()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeOperationLogListFailed, err.Error())
	}

	// 读取所有日志记录
	logs, err := s.readAllLogs(logFiles)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeOperationLogListFailed, err.Error())
	}

	// 按时间排序（降序）
	sort.Slice(logs, func(i, j int) bool {
		return logs[i].OperationTime.Timestamp() > logs[j].OperationTime.Timestamp()
	})

	// 过滤日志
	filteredLogs := s.filterLogs(logs, req)

	// 计算总数
	res.Total = len(filteredLogs)

	// 分页处理
	startIndex := (req.Page - 1) * req.PageSize
	endIndex := startIndex + req.PageSize
	if startIndex >= len(filteredLogs) {
		return res, nil
	}
	if endIndex > len(filteredLogs) {
		endIndex = len(filteredLogs)
	}

	// 设置结果
	for _, log := range filteredLogs[startIndex:endIndex] {
		resultText := "成功"
		if log.OperationResult == 0 {
			resultText = "失败"
		}

		res.List = append(res.List, v1.OperationLogListItem{
			Id:              log.Id,
			UserId:          log.UserId,
			Username:        log.Username,
			OperationIp:     log.OperationIp,
			OperationTime:   log.OperationTime,
			Module:          log.Module,
			Action:          log.Action,
			OperationResult: log.OperationResult,
			ResultText:      resultText,
		})
	}

	return res, nil
}

// Delete 删除操作日志
func (s *logImpl) Delete(ctx context.Context, req *v1.OperationLogDeleteReq) (res *v1.OperationLogDeleteRes, err error) {
	res = &v1.OperationLogDeleteRes{}

	// 获取所有日志文件
	logFiles, err := s.getLogFiles()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeOperationLogDeleteFailed, err.Error())
	}

	// 读取所有日志记录
	logs, err := s.readAllLogs(logFiles)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeOperationLogDeleteFailed, err.Error())
	}

	// 检查日志是否存在
	found := false
	for _, log := range logs {
		if log.Id == req.Id {
			found = true
			break
		}
	}

	if !found {
		return nil, gerror.NewCode(consts.CodeOperationLogDeleteFailed, "操作日志不存在")
	}

	// 从所有日志中排除要删除的日志
	var newLogs []logRecord
	for _, log := range logs {
		if log.Id != req.Id {
			newLogs = append(newLogs, log)
		}
	}

	// 清空日志文件并重新写入
	err = s.rewriteAllLogs(newLogs)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeOperationLogDeleteFailed, err.Error())
	}

	return res, nil
}

// Export 导出操作日志
func (s *logImpl) Export(ctx context.Context, req *v1.OperationLogExportReq) (res *v1.OperationLogExportRes, err error) {
	res = &v1.OperationLogExportRes{}

	// 获取所有日志文件
	logFiles, err := s.getLogFiles()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeOperationLogExportFailed, err.Error())
	}

	// 读取所有日志记录
	logs, err := s.readAllLogs(logFiles)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeOperationLogExportFailed, err.Error())
	}

	// 过滤日志
	filteredLogs := s.filterLogsForExport(logs, req)

	// 确保导出目录存在
	exportDir := "resource/export"
	if !gfile.Exists(exportDir) {
		if err = gfile.Mkdir(exportDir); err != nil {
			return nil, gerror.NewCode(consts.CodeOperationLogExportFailed, err.Error())
		}
	}

	// 生成CSV文件
	filename := fmt.Sprintf("operation_log_%s.csv", gtime.Now().Format("YmdHis"))
	filepath := filepath.Join(exportDir, filename)

	// 创建文件
	file, err := os.Create(filepath)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeOperationLogExportFailed, err.Error())
	}
	defer file.Close()

	// 写入BOM，解决Excel打开中文乱码问题
	_, err = file.Write([]byte{0xEF, 0xBB, 0xBF})
	if err != nil {
		return nil, gerror.NewCode(consts.CodeOperationLogExportFailed, err.Error())
	}

	// 创建CSV写入器
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入CSV头
	headers := []string{"ID", "用户ID", "操作人", "操作IP", "操作时间", "操作模块", "操作类型", "操作结果"}
	if err = writer.Write(headers); err != nil {
		return nil, gerror.NewCode(consts.CodeOperationLogExportFailed, err.Error())
	}

	// 写入数据行
	for _, item := range filteredLogs {
		resultText := "成功"
		if item.OperationResult == 0 {
			resultText = "失败"
		}

		row := []string{
			gconv.String(item.Id),
			gconv.String(item.UserId),
			item.Username,
			item.OperationIp,
			item.OperationTime.String(),
			item.Module,
			item.Action,
			resultText,
		}

		if err = writer.Write(row); err != nil {
			return nil, gerror.NewCode(consts.CodeOperationLogExportFailed, err.Error())
		}
	}

	// 设置下载URL
	res.Url = "/export/" + filename
	return res, nil
}

// Record 记录操作日志
func (s *logImpl) Record(ctx context.Context, userId int, username, module, action string, result int, details string) error {
	// 获取客户端IP
	ip := ""

	// 尝试从上下文中获取请求对象
	if reqObj := ctx.Value(middleware.RequestObjectKey); reqObj != nil {
		// 尝试类型断言为http请求对象
		if request, ok := reqObj.(*ghttp.Request); ok {
			ip = request.GetClientIp()
		}
	}

	// 如果上面的方式无法获取IP，再尝试从标准方式获取
	if ip == "" {
		request := g.RequestFromCtx(ctx)
		if request != nil {
			ip = request.GetClientIp()
		}
	}

	// 创建日志记录
	record := logRecord{
		Id:              s.generateId(),
		UserId:          userId,
		Username:        username,
		OperationIp:     ip, // 使用从请求中获取的IP
		OperationTime:   gtime.Now(),
		Module:          module,
		Action:          action,
		OperationResult: result,
	}

	// 异步写入日志，不阻塞主业务流程
	go func(logData logRecord) {
		if err := s.writeLog(logData); err != nil {
			g.Log().Error(context.Background(), "记录操作日志失败:", err)
		}
	}(record)

	return nil
}

// logRecord 日志记录结构
type logRecord struct {
	Id              int         `json:"id"`
	UserId          int         `json:"userId"`
	Username        string      `json:"username"`
	OperationIp     string      `json:"operationIp"`
	OperationTime   *gtime.Time `json:"operationTime"`
	Module          string      `json:"module"`
	Action          string      `json:"action"`
	OperationResult int         `json:"operationResult"`
}

// generateId 生成唯一ID
func (s *logImpl) generateId() int {
	return int(gtime.TimestampMilli())
}

// getLogFiles 获取所有日志文件路径
func (s *logImpl) getLogFiles() ([]string, error) {
	pattern := filepath.Join(logDir, "operation_*.log")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	// 按文件修改时间排序（最新的在前）
	sort.Slice(files, func(i, j int) bool {
		infoI, _ := os.Stat(files[i])
		infoJ, _ := os.Stat(files[j])
		return infoI.ModTime().After(infoJ.ModTime())
	})

	return files, nil
}

// getCurrentLogFile 获取当前要写入的日志文件
func (s *logImpl) getCurrentLogFile() (string, error) {
	// 日志文件命名格式：operation_YYYYMMDD.log
	today := gtime.Now().Format("Ymd")
	filename := filepath.Join(logDir, fmt.Sprintf("operation_%s.log", today))

	// 检查文件是否存在，不存在则创建
	if !gfile.Exists(filename) {
		// 创建空文件
		if err := os.WriteFile(filename, []byte{}, 0666); err != nil {
			return "", err
		}
		return filename, nil
	}

	// 检查文件行数，超过最大行数则创建新文件
	content := gfile.GetContents(filename)
	lineCount := len(strings.Split(content, "\n"))
	if lineCount < maxLogLines {
		return filename, nil
	}

	// 行数超过最大值，创建新文件（加上时间戳以区分）
	timestamp := gtime.Now().Format("YmdHis")
	newFilename := filepath.Join(logDir, fmt.Sprintf("operation_%s_%s.log", today, timestamp))
	// 创建空文件
	if err := os.WriteFile(newFilename, []byte{}, 0666); err != nil {
		return "", err
	}

	return newFilename, nil
}

// cleanOldLogs 清理旧日志文件
func (s *logImpl) cleanOldLogs() error {
	files, err := s.getLogFiles()
	if err != nil {
		return err
	}

	// 如果文件数量超过最大值，删除最旧的文件
	if len(files) > maxLogFiles {
		// 文件已按时间排序，删除后面的旧文件
		for i := maxLogFiles; i < len(files); i++ {
			if err := gfile.Remove(files[i]); err != nil {
				return err
			}
		}
	}

	return nil
}

// writeLog 写入单条日志
func (s *logImpl) writeLog(record logRecord) error {
	logMutex.Lock()
	defer logMutex.Unlock()

	// 清理旧日志
	if err := s.cleanOldLogs(); err != nil {
		return err
	}

	// 获取当前日志文件
	filename, err := s.getCurrentLogFile()
	if err != nil {
		return err
	}

	// 将记录转换为JSON
	jsonBytes, err := json.Marshal(record)
	if err != nil {
		return err
	}

	// 追加到文件
	if err := gfile.PutContentsAppend(filename, string(jsonBytes)+"\n"); err != nil {
		return err
	}

	return nil
}

// readAllLogs 读取所有日志
func (s *logImpl) readAllLogs(files []string) ([]logRecord, error) {
	var allLogs []logRecord

	for _, file := range files {
		content := gfile.GetContents(file)
		if content == "" {
			continue
		}

		// 按行读取
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			if line = strings.TrimSpace(line); line == "" {
				continue
			}

			var record logRecord
			if err := json.Unmarshal([]byte(line), &record); err != nil {
				g.Log().Warning(context.Background(), "解析日志行失败:", err, "行内容:", line)
				continue
			}

			allLogs = append(allLogs, record)
		}
	}

	return allLogs, nil
}

// filterLogs 过滤日志记录
func (s *logImpl) filterLogs(logs []logRecord, req *v1.OperationLogListReq) []logRecord {
	var filtered []logRecord

	for _, log := range logs {
		// 用户名过滤
		if req.Username != "" && !gstr.Contains(log.Username, req.Username) {
			continue
		}

		// 模块过滤
		if req.Module != "" && log.Module != req.Module {
			continue
		}

		// 操作类型过滤
		if req.Action != "" && log.Action != req.Action {
			continue
		}

		// 关键字过滤（用户名、操作类型）
		if req.Keyword != "" && !gstr.Contains(log.Username, req.Keyword) && !gstr.Contains(log.Action, req.Keyword) {
			continue
		}

		// 结果过滤
		if req.Result != "" {
			reqResult := gconv.Int(req.Result)
			if log.OperationResult != reqResult {
				continue
			}
		}

		// 时间范围过滤
		if req.StartTime != "" {
			startTime := gtime.NewFromStr(req.StartTime)
			if startTime != nil && log.OperationTime.Before(startTime) {
				continue
			}
		}

		if req.EndTime != "" {
			endTime := gtime.NewFromStr(req.EndTime)
			if endTime != nil && log.OperationTime.After(endTime) {
				continue
			}
		}

		filtered = append(filtered, log)
	}

	return filtered
}

// filterLogsForExport 过滤用于导出的日志记录
func (s *logImpl) filterLogsForExport(logs []logRecord, req *v1.OperationLogExportReq) []logRecord {
	var filtered []logRecord

	for _, log := range logs {
		// 用户名过滤
		if req.Username != "" && !gstr.Contains(log.Username, req.Username) {
			continue
		}

		// 模块过滤
		if req.Module != "" && log.Module != req.Module {
			continue
		}

		// 操作类型过滤
		if req.Action != "" && log.Action != req.Action {
			continue
		}

		// 关键字过滤（用户名、操作类型）
		if req.Keyword != "" && !gstr.Contains(log.Username, req.Keyword) && !gstr.Contains(log.Action, req.Keyword) {
			continue
		}

		// 结果过滤
		if req.Result != "" {
			reqResult := gconv.Int(req.Result)
			if log.OperationResult != reqResult {
				continue
			}
		}

		// 时间范围过滤
		if req.StartTime != "" {
			startTime := gtime.NewFromStr(req.StartTime)
			if startTime != nil && log.OperationTime.Before(startTime) {
				continue
			}
		}

		if req.EndTime != "" {
			endTime := gtime.NewFromStr(req.EndTime)
			if endTime != nil && log.OperationTime.After(endTime) {
				continue
			}
		}

		filtered = append(filtered, log)
	}

	return filtered
}

// rewriteAllLogs 重写所有日志文件
func (s *logImpl) rewriteAllLogs(logs []logRecord) error {
	logMutex.Lock()
	defer logMutex.Unlock()

	// 删除所有现有日志文件
	files, err := s.getLogFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := gfile.Remove(file); err != nil {
			return err
		}
	}

	// 如果没有日志需要写入，直接返回
	if len(logs) == 0 {
		return nil
	}

	// 按时间排序（降序）
	sort.Slice(logs, func(i, j int) bool {
		return logs[i].OperationTime.Timestamp() > logs[j].OperationTime.Timestamp()
	})

	// 重新写入日志，按日期分组
	dateGroups := make(map[string][]logRecord)
	for _, log := range logs {
		date := log.OperationTime.Format("Ymd")
		dateGroups[date] = append(dateGroups[date], log)
	}

	// 写入每个日期组
	for date, group := range dateGroups {
		// 如果一个日期的日志超过最大行数，分多个文件
		chunks := s.chunkSlice(group, maxLogLines)

		for i, chunk := range chunks {
			var filename string
			if i == 0 {
				filename = filepath.Join(logDir, fmt.Sprintf("operation_%s.log", date))
			} else {
				filename = filepath.Join(logDir, fmt.Sprintf("operation_%s_%d.log", date, i))
			}

			// 创建文件
			if err := os.WriteFile(filename, []byte{}, 0666); err != nil {
				return err
			}

			// 写入每条日志
			var content strings.Builder
			for _, record := range chunk {
				jsonBytes, err := json.Marshal(record)
				if err != nil {
					return err
				}
				content.WriteString(string(jsonBytes) + "\n")
			}

			if err := gfile.PutContents(filename, content.String()); err != nil {
				return err
			}
		}
	}

	return nil
}

// chunkSlice 将切片分块
func (s *logImpl) chunkSlice(slice []logRecord, chunkSize int) [][]logRecord {
	var chunks [][]logRecord
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}
