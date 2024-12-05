<script setup lang="ts">
import { ref, onBeforeUnmount } from 'vue'
import { ElMessage } from 'element-plus'
import 'element-plus/dist/index.css'
import { AddTextWatermark, AddImageWatermark, UploadImage, GetImagePreview } from '../wailsjs/go/main/App'
import { watermark, color } from '../wailsjs/go/models'
import type { Ref } from 'vue'

// 水印类型
const watermarkType = ref<'text' | 'image'>('text')

// 水印设置
const watermarkSettings = ref({
  // 文本水印设置
  text: '',
  textSize: 40,
  textColor: '#000000',
  fontFamily: 'Arial',
  
  // 图片水印设置
  imageFile: null as File | null,
  scale: 1.0,
  
  // 通用设置
  opacity: 0.5,
  angle: 45,
  spacing: 100,
  position: 'center' as 'center' | 'topLeft' | 'topRight' | 'bottomLeft' | 'bottomRight' | 'tiled',
  margin: 20,
})

// 水印位置选项
const positionOptions = [
  { label: '居中', value: 'center' },
  { label: '左上角', value: 'topLeft' },
  { label: '右上角', value: 'topRight' },
  { label: '左下角', value: 'bottomLeft' },
  { label: '右下角', value: 'bottomRight' },
  { label: '平铺', value: 'tiled' },
] as const

// 选中的图片文件
const selectedFiles = ref<{
  name: string;
  path: string;
  size: number;
  extension: string;
}[]>([])

// 预览图片URL
const previewUrl = ref<string>('')

// 文件输入引用
const fileInput = ref<HTMLInputElement | null>(null)

// 点击选择文件
const handleClickFileInput = () => {
  if (fileInput.value) {
    fileInput.value.click()
  }
}

// 处理图片拖拽
const handleDrop = async (e: DragEvent) => {
  e.preventDefault()
  const files = e.dataTransfer?.files
  if (files) {
    await handleFiles(Array.from(files))
  }
}

// 处理文件选择
const handleFileSelect = async (e: Event) => {
  const files = (e.target as HTMLInputElement).files
  if (files) {
    await handleFiles(Array.from(files))
  }
}

// 处理文件
const handleFiles = async (files: File[]) => {
  const imageFiles = files.filter(file => 
    file.type.startsWith('image/')
  )
  
  if (imageFiles.length === 0) {
    ElMessage.error('请选择图片文件')
    return
  }

  try {
    const uploadedFiles = []
    for (const file of imageFiles) {
      const arrayBuffer = await file.arrayBuffer()
      const uint8Array = new Uint8Array(arrayBuffer)
      const result = await UploadImage(Array.from(uint8Array), file.name)
      if (result) {
        uploadedFiles.push(result)
      }
    }

    selectedFiles.value = uploadedFiles
    // 预览第一张图片
    if (uploadedFiles[0]) {
      try {
        // Get base64 preview from backend
        const preview = await GetImagePreview(uploadedFiles[0].path)
        previewUrl.value = preview
      } catch (error) {
        console.error('Failed to get image preview:', error)
        ElMessage.error('Failed to load image preview')
      }
    }
  } catch (error) {
    ElMessage.error('上传图片失败: ' + error)
  }
}

// 处理水印图片选择
const handleWatermarkImageSelect = async (uploadFile: any) => {
  if (uploadFile.raw) {
    watermarkSettings.value.imageFile = uploadFile.raw
    try {
      // Get base64 preview from backend
      const preview = await GetImagePreview(uploadFile.raw.path)
      previewUrl.value = preview
    } catch (error) {
      console.error('Failed to get image preview:', error)
      ElMessage.error('Failed to load image preview')
    }
  }
}

// 格式化工具提示
const formatTooltip = (value: number): string => `${value}px`

// 应用水印
const applyWatermark = async () => {
  if (selectedFiles.value.length === 0) {
    ElMessage.error('请先选择要处理的图片')
    return
  }

  try {
    for (const file of selectedFiles.value) {
      if (watermarkType.value === 'image' && watermarkSettings.value.imageFile) {
        // 应用图片水印
        const options = new watermark.WatermarkOptions({
          scale: watermarkSettings.value.scale,
          opacity: watermarkSettings.value.opacity,
          angle: watermarkSettings.value.angle,
          spacing: watermarkSettings.value.spacing,
          position: watermarkSettings.value.position,
          margin: watermarkSettings.value.margin,
        })
        await AddImageWatermark(file.path, watermarkSettings.value.imageFile.name, options)
      } else {
        // 应用文本水印
        const options = new watermark.WatermarkOptions({
          text: watermarkSettings.value.text,
          textSize: watermarkSettings.value.textSize,
          textColor: new color.RGBA({
            R: parseInt(watermarkSettings.value.textColor.slice(1, 3), 16),
            G: parseInt(watermarkSettings.value.textColor.slice(3, 5), 16),
            B: parseInt(watermarkSettings.value.textColor.slice(5, 7), 16),
            A: 255
          }),
          fontFamily: watermarkSettings.value.fontFamily,
          opacity: watermarkSettings.value.opacity,
          angle: watermarkSettings.value.angle,
          spacing: watermarkSettings.value.spacing,
          position: watermarkSettings.value.position,
          margin: watermarkSettings.value.margin,
        })
        await AddTextWatermark(file.path, options)
      }
    }
    ElMessage.success('水印添加成功')
  } catch (error) {
    ElMessage.error('水印添加失败: ' + error)
  }
}

onBeforeUnmount(() => {
  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value)
  }
})
</script>

<template>
  <el-container class="app-container">
    <el-header>
      <h1>图片水印工具</h1>
    </el-header>
    
    <el-main>
      <el-row :gutter="20">
        <!-- 左侧：图片上传和预览区域 -->
        <el-col :span="16">
          <el-card class="upload-area" 
                  @drop.prevent="handleDrop" 
                  @dragover.prevent
                  @dragleave.prevent>
            <div v-if="!previewUrl" class="upload-placeholder">
              <el-icon><upload-filled /></el-icon>
              <div>拖拽图片到此处或点击上传</div>
              <input type="file" 
                     accept="image/*" 
                     multiple 
                     @change="handleFileSelect" 
                     style="display: none" 
                     ref="fileInput">
              <el-button type="primary" @click="handleClickFileInput">
                选择图片
              </el-button>
            </div>
            <div v-else class="preview-container">
              <img :src="previewUrl" class="preview-image">
              <div class="preview-info">
                <p>已选择 {{ selectedFiles.length }} 个文件</p>
                <el-button type="primary" link @click="handleClickFileInput">
                  重新选择
                </el-button>
              </div>
            </div>
          </el-card>
        </el-col>

        <!-- 右侧：水印设置面板 -->
        <el-col :span="8">
          <el-card class="settings-panel">
            <el-form label-position="top">
              <el-tabs v-model="watermarkType">
                <el-tab-pane label="文本水印" name="text">
                  <el-form-item label="水印文本">
                    <el-input v-model="watermarkSettings.text" />
                  </el-form-item>
                  
                  <el-form-item label="字体大小">
                    <el-slider v-model="watermarkSettings.textSize" 
                              :min="12" 
                              :max="100"
                              :format-tooltip="formatTooltip" />
                  </el-form-item>
                  
                  <el-form-item label="文字颜色">
                    <el-color-picker v-model="watermarkSettings.textColor" 
                                    show-alpha />
                  </el-form-item>
                  
                  <el-form-item label="字体">
                    <el-select v-model="watermarkSettings.fontFamily">
                      <el-option label="Arial" value="Arial" />
                      <el-option label="Times New Roman" value="Times New Roman" />
                      <el-option label="Courier New" value="Courier New" />
                    </el-select>
                  </el-form-item>
                </el-tab-pane>
                
                <el-tab-pane label="图片水印" name="image">
                  <el-form-item label="水印图片">
                    <el-upload
                      action=""
                      :auto-upload="false"
                      :show-file-list="true"
                      accept="image/*"
                      :limit="1"
                      :on-change="handleWatermarkImageSelect">
                      <el-button type="primary">选择图片</el-button>
                      <template #tip>
                        <div class="el-upload__tip">支持 jpg/png 格式图片</div>
                      </template>
                    </el-upload>
                    <div v-if="previewUrl" class="preview-image">
                      <img :src="previewUrl" alt="水印图片预览" style="max-width: 200px; margin-top: 10px;" />
                    </div>
                  </el-form-item>
                  
                  <el-form-item label="缩放比例">
                    <el-slider v-model="watermarkSettings.scale" 
                              :min="0.1" 
                              :max="2"
                              :step="0.1" />
                  </el-form-item>
                </el-tab-pane>
              </el-tabs>
              
              <el-divider>通用设置</el-divider>
              
              <el-form-item label="透明度">
                <el-slider v-model="watermarkSettings.opacity" 
                          :min="0" 
                          :max="1"
                          :step="0.1" />
              </el-form-item>
              
              <el-form-item label="旋转角度">
                <el-slider v-model="watermarkSettings.angle" 
                          :min="-180" 
                          :max="180" />
              </el-form-item>
              
              <el-form-item label="水印位置">
                <el-select v-model="watermarkSettings.position">
                  <el-option
                    v-for="option in positionOptions"
                    :key="option.value"
                    :label="option.label"
                    :value="option.value" />
                </el-select>
              </el-form-item>
              
              <el-form-item label="水印间距">
                <el-slider v-model="watermarkSettings.spacing" 
                          :min="50" 
                          :max="300"
                          :format-tooltip="formatTooltip" />
              </el-form-item>

              <el-form-item label="边距">
                <el-slider v-model="watermarkSettings.margin" 
                          :min="0" 
                          :max="100"
                          :format-tooltip="formatTooltip" />
              </el-form-item>

              <el-form-item>
                <el-button type="primary" 
                           @click="applyWatermark" 
                           :disabled="selectedFiles.length === 0 || 
                                     (watermarkType === 'text' && !watermarkSettings.text) ||
                                     (watermarkType === 'image' && !watermarkSettings.imageFile)">
                  应用水印
                </el-button>
              </el-form-item>
            </el-form>
          </el-card>
        </el-col>
      </el-row>
    </el-main>
  </el-container>
</template>

<style scoped>
.app-container {
  height: 100vh;
  background-color: #f5f7fa;
}

.el-header {
  background-color: white;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.12);
  margin-bottom: 20px;
  display: flex;
  align-items: center;
}

.el-main {
  padding: 20px;
}

.upload-area {
  min-height: 400px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.upload-placeholder {
  text-align: center;
  color: #909399;
}

.upload-placeholder .el-icon {
  font-size: 48px;
  margin-bottom: 10px;
}

.preview-container {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.preview-image {
  max-width: 100%;
  max-height: 300px;
  object-fit: contain;
  margin-bottom: 20px;
}

.preview-info {
  text-align: center;
}

.settings-panel {
  height: 100%;
}

.el-divider {
  margin: 24px 0;
}
</style>
