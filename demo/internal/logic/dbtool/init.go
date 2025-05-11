package dbtool

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/glog"
)

// InitDatabase 初始化数据库
// 在应用启动时检查数据库表是否存在或为空，如果是则初始化数据库
func InitDatabase(ctx context.Context) error {
	// 检查用户表是否存在数据
	count, err := g.DB().Ctx(ctx).Model("user").Count()
	if err != nil {
		if strings.Contains(err.Error(), "doesn't exist") {
			// 表不存在，需要初始化
			glog.Info(ctx, "数据库表不存在，开始初始化数据库...")
			return executeSqlFile(ctx)
		}
		// 其他错误
		return err
	}

	if count == 0 {
		// 表存在但没有数据，需要初始化
		glog.Info(ctx, "数据库表为空，开始初始化数据...")
		return executeSqlFile(ctx)
	}

	glog.Info(ctx, "数据库已初始化，用户数据记录数:", count)

	// 检查browse_history表是否存在，不存在则创建
	err = checkAndCreateBrowseHistoryTable(ctx)
	if err != nil {
		return err
	}

	return nil
}

// executeSqlFile 执行SQL初始化文件
func executeSqlFile(ctx context.Context) error {
	// SQL文件路径
	sqlPath := "manifest/sql/init.sql"

	// 确保文件存在
	if !gfile.Exists(sqlPath) {
		// 尝试查找上级目录
		dir, _ := os.Getwd()
		parentDir := filepath.Dir(dir)
		alternativePath := filepath.Join(parentDir, sqlPath)

		if gfile.Exists(alternativePath) {
			sqlPath = alternativePath
		} else {
			return gerror.New("初始化SQL文件不存在:" + sqlPath)
		}
	}

	// 读取SQL文件内容
	sqlContent := gfile.GetContents(sqlPath)
	if sqlContent == "" {
		return gerror.New("SQL文件为空:" + sqlPath)
	}

	// 分割SQL语句
	sqlStatements := splitSql(sqlContent)

	// 开始事务
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, sqlStmt := range sqlStatements {
			if strings.TrimSpace(sqlStmt) == "" {
				continue
			}

			_, err := tx.Exec(sqlStmt)
			if err != nil {
				// 忽略已存在表的错误
				if strings.Contains(err.Error(), "already exists") ||
					strings.Contains(err.Error(), "Duplicate entry") {
					glog.Warning(ctx, "SQL执行警告:", err.Error())
					continue
				}
				glog.Error(ctx, "执行SQL失败:", err.Error(), "SQL:", sqlStmt)
				return err
			}
		}
		return nil
	})

	if err != nil {
		glog.Error(ctx, "数据库初始化失败:", err.Error())
		return err
	}

	glog.Info(ctx, "数据库初始化成功!")
	return nil
}

// splitSql 分割SQL语句为单独的命令
func splitSql(sql string) []string {
	// 替换SQL注释
	lines := strings.Split(sql, "\n")
	var filteredLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}
		filteredLines = append(filteredLines, line)
	}

	sql = strings.Join(filteredLines, "\n")

	// 按分号分割SQL语句
	statements := strings.Split(sql, ";")
	var result []string

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" {
			result = append(result, stmt+";")
		}
	}

	return result
}

// 检查并创建浏览历史记录表
func checkAndCreateBrowseHistoryTable(ctx context.Context) error {
	// 检查表是否存在
	sql := "SHOW TABLES LIKE 'browse_history'"
	result, err := g.DB().GetAll(ctx, sql)
	if err != nil {
		return err
	}

	// 表已存在，无需创建
	if len(result) > 0 {
		return nil
	}

	// 创建表
	sql = `
	CREATE TABLE browse_history (
		id INT NOT NULL AUTO_INCREMENT COMMENT '记录ID',
		client_id INT NOT NULL COMMENT '客户ID',
		content_id INT NOT NULL COMMENT '内容ID',
		content_type VARCHAR(20) NOT NULL COMMENT '内容类型',
		browse_time DATETIME NOT NULL COMMENT '浏览时间',
		PRIMARY KEY (id),
		INDEX idx_client_time (client_id, browse_time DESC),
		INDEX idx_content (content_id, content_type)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='浏览历史记录表';
	`

	_, err = g.DB().Exec(ctx, sql)
	if err != nil {
		return err
	}

	// 记录日志
	g.Log().Info(ctx, "浏览历史记录表 browse_history 创建成功")
	return nil
}
