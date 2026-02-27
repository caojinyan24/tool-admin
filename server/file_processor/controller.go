package file_processor

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"tooladmin/server/common/logger"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

var FileProcessor = NewFileProcessorController()

// FileProcessorController 处理文件处理相关请求
type FileProcessorController struct{}

// NewFileProcessorController 创建新的文件处理器控制器
func NewFileProcessorController() *FileProcessorController {
	return &FileProcessorController{}
}

// MergePDF 处理PDF合并请求
func (*FileProcessorController) MergePDF(ctx context.Context, c *app.RequestContext) {
	logger.Info(ctx, "合并PDF文件请求")
	// 解析多部分表单
	form, err := c.MultipartForm()
	if err != nil {
		logger.Error(ctx, "解析文件失败", err)
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "无法解析文件",
		})
		return
	}

	// 获取上传的文件
	files := form.File["files"]
	if len(files) < 2 {
		logger.Error(ctx, "上传文件数量不足", err)
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "至少需要上传2个PDF文件",
		})
		return
	}

	// 创建临时目录存储上传的文件
	tempDir := fmt.Sprintf("./temp/%d", time.Now().UnixNano())
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		logger.Error(ctx, "创建临时目录失败", err)
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "无法创建临时目录",
		})
		return
	}
	defer os.RemoveAll(tempDir) // 处理完成后清理临时文件

	// 保存上传的文件
	var inFiles []string
	for i, file := range files {
		// 检查文件类型
		if filepath.Ext(file.Filename) != ".pdf" {
			logger.Error(ctx, "文件类型错误", err)
			c.JSON(http.StatusBadRequest, map[string]string{
				"error": fmt.Sprintf("文件 %s 不是PDF文件", file.Filename),
			})
			return
		}

		// 保存文件
		tempPath := filepath.Join(tempDir, fmt.Sprintf("file%d.pdf", i))
		src, err := file.Open()
		if err != nil {
			logger.Error(ctx, "打开上传的文件失败", err)
			c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "无法打开上传的文件",
			})
			return
		}
		defer src.Close()

		dst, err := os.Create(tempPath)
		if err != nil {
			logger.Error(ctx, "创建临时文件失败", err)
			c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "无法创建临时文件",
			})
			return
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			logger.Error(ctx, "复制文件失败", err)
			c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "无法保存上传的文件",
			})
			return
		}

		inFiles = append(inFiles, tempPath)
	}

	// 合并PDF文件
	outFile := filepath.Join(tempDir, "merged.pdf")
	if err := api.MergeCreateFile(inFiles, outFile, false, nil); err != nil {
		logger.Error(ctx, "合并PDF失败", err)
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("合并PDF失败: %v", err),
		})
		return
	}

	// 读取合并后的文件
	mergedBytes, err := os.ReadFile(outFile)
	if err != nil {
		logger.Error(ctx, "读取合并后的文件失败", err)
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "无法读取合并后的文件",
		})
		return
	}

	// 设置响应头
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=merged.pdf")
	c.Header("Content-Length", fmt.Sprintf("%d", len(mergedBytes)))

	// 返回合并后的文件
	c.Write(mergedBytes)
}
