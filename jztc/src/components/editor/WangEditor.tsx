import React, { useState, useEffect } from 'react';
import { Editor, Toolbar } from '@wangeditor/editor-for-react';
import { IDomEditor, IEditorConfig, IToolbarConfig } from '@wangeditor/editor';
import '@wangeditor/editor/dist/css/style.css';
import { message } from 'antd';
import FileSelector from '@/components/FileSelector';
import './WangEditor.css';

interface WangEditorProps {
  value?: string;
  onChange?: (html: string) => void;
  height?: number;
  placeholder?: string;
}

const WangEditor: React.FC<WangEditorProps> = ({
  value = '',
  onChange,
  height = 500,
  placeholder = '请输入内容...'
}) => {
  // 编辑器实例，必须用 useState 存储
  const [editor, setEditor] = useState<IDomEditor | null>(null);
  // 编辑器内容
  const [html, setHtml] = useState<string>(value);
  // 文件选择弹窗
  const [fileModalVisible, setFileModalVisible] = useState<boolean>(false);
  // 当前插入函数
  const [currentInsertFn, setCurrentInsertFn] = useState<((url: string, alt: string, href: string) => void) | null>(null);

  // 监听 value 变化
  useEffect(() => {
    if (value !== html) {
      setHtml(value);
    }
  }, [value]);

  // 打开文件选择弹窗
  const openFileModal = (insertFn: (url: string, alt: string, href: string) => void) => {
    if (typeof insertFn === 'function') {
      setCurrentInsertFn(() => insertFn); // 使用函数形式设置状态，避免引用问题
      setFileModalVisible(true);
    } else {
      console.error('insertFn is not a function', insertFn);
      message.error('插入函数无效，请重试');
    }
  };

  // 关闭文件选择弹窗
  const closeFileModal = () => {
    setFileModalVisible(false);
    setCurrentInsertFn(null);
  };

  // 处理文件选择
  const handleSelectFile = (url: string) => {
    if (currentInsertFn && typeof currentInsertFn === 'function') {
      try {
        // 只传入url，不传入alt和href
        currentInsertFn(url, '', '');
      } catch (error) {
        console.error('插入图片失败:', error);
        message.error('插入图片失败，请重试');
      }
    } else {
      console.error('currentInsertFn is not a function', currentInsertFn);
      message.error('插入函数无效，请关闭对话框重试');
    }
  };

  // 工具栏配置
  const toolbarConfig: Partial<IToolbarConfig> = {
    toolbarKeys: [
      'headerSelect',
      'blockquote',
      '|',
      'bold',
      'italic',
      'underline',
      'through',
      'color',
      'bgColor',
      '|',
      'bulletedList',
      'numberedList',
      'todo',
      '|',
      'justifyLeft',
      'justifyCenter',
      'justifyRight',
      '|',
      {
        key: 'group-image', // 必填，要以 group 开头
        title: '图片', // 必填
        menuKeys: ['uploadImage'] // 只保留自定义的上传图片功能
      },
      'insertTable',
      '|',
      'code',
      'codeBlock',
      '|',
      'undo',
      'redo',
    ],
    excludeKeys: [
      'insertImage' // 排除默认的插入网络图片按钮
    ]
  };

  // 编辑器配置
  const editorConfig: Partial<IEditorConfig> = {
    placeholder: placeholder,
    autoFocus: false,
    MENU_CONF: {
      uploadImage: {
        // 使用自定义选择图片功能
        customBrowseAndUpload(insertFn: (url: string, alt: string, href: string) => void) {
          // 打开文件选择弹窗
          openFileModal(insertFn);
        }
      },
      insertImage: {
        // 自定义插入网络图片
        customBrowseAndUpload(insertFn: (url: string, alt: string, href: string) => void) {
          // 也使用文件选择弹窗替代默认的插入网络图片功能
          openFileModal(insertFn);
        }
      }
    }
  };

  // 处理内容变化
  const handleChange = (editor: IDomEditor) => {
    // 获取原始HTML
    let newHtml = editor.getHtml();
    
    // 简化图片标签 - 将复杂的img标签替换为只有src的简单版本
    newHtml = newHtml.replace(/<img[^>]+src="([^"]+)"[^>]*>/g, '<img src="$1">');
    
    setHtml(newHtml);
    if (onChange) {
      onChange(newHtml);
    }
  };

  // 组件销毁时，销毁编辑器
  useEffect(() => {
    return () => {
      if (editor == null) return;
      editor.destroy();
      setEditor(null);
    };
  }, [editor]);

  return (
    <>
      <div style={{ border: '1px solid #ccc', zIndex: 100 }}>
        <Toolbar
          editor={editor}
          defaultConfig={toolbarConfig}
          mode="default"
          style={{ borderBottom: '1px solid #ccc' }}
        />
        <Editor
          defaultConfig={editorConfig}
          value={html}
          onCreated={setEditor}
          onChange={handleChange}
          mode="default"
          style={{ height: `${height}px`, overflowY: 'hidden' }}
        />
      </div>

      {/* 使用通用文件选择器组件 */}
      <FileSelector 
        visible={fileModalVisible}
        onCancel={closeFileModal}
        onSelect={handleSelectFile}
        title="选择图片"
        accept="image/*"
      />
    </>
  );
};

export default WangEditor; 