<template>
  <div class="container">
    <div class="header">
      <div class="file-inputs">
        <div class="file-input">
          <label>Left File:</label>
          <input type="file" @change="handleLeftFileChange" />
        </div>
        <div class="file-input">
          <label>Right File:</label>
          <input type="file" @change="handleRightFileChange" />
        </div>
      </div>
      <div class="actions">
        <button @click="compareFiles" :disabled="!canCompare">Compare Files</button>
        <button @click="applyLeftChanges" :disabled="!diffResult">Apply Left Changes</button>
        <button @click="applyRightChanges" :disabled="!diffResult">Apply Right Changes</button>
        <button @click="mergeFiles" :disabled="!diffResult">Merge Files</button>
      </div>
    </div>

    <div class="diff-container" v-if="diffResult">
      <div class="diff-pane left">
        <div
          v-for="line in diffResult.leftLines"
          :key="line.lineNumber"
          :class="['diff-line', line.type]"
        >
          <span class="line-number">{{ line.lineNumber }}</span>
          <span class="line-content">{{ line.content }}</span>
        </div>
      </div>
      <div class="diff-pane right">
        <div
          v-for="line in diffResult.rightLines"
          :key="line.lineNumber"
          :class="['diff-line', line.type]"
        >
          <span class="line-number">{{ line.lineNumber }}</span>
          <span class="line-content">{{ line.content }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { CompareFiles } from './api'
import { ApplyChanges, MergeFiles } from '../wailsjs/go/main/App';

type DiffLineType = 'same' | 'added' | 'removed'

interface DiffLine {
  lineNumber: number
  content: string
  type: DiffLineType
}

interface DiffResult {
  leftLines: DiffLine[]
  rightLines: DiffLine[]
}

const leftFilePath = ref<string>('')
const rightFilePath = ref<string>('')
const diffResult = ref<DiffResult | null>(null)

const canCompare = computed(() => leftFilePath.value && rightFilePath.value)

const handleLeftFileChange = (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (file) {
    leftFilePath.value = file.name
  }
}

const handleRightFileChange = (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (file) {
    rightFilePath.value = file.name
  }
}

const compareFiles = async () => {
  try {
    diffResult.value = await CompareFiles(leftFilePath.value, rightFilePath.value)
  } catch (error) {
    console.error('Error comparing files:', error)
    alert('Error comparing files')
  }
}

const applyLeftChanges = async () => {
  try {
    await ApplyChanges(rightFilePath.value, leftFilePath.value, 'left')
    await compareFiles()
  } catch (error) {
    console.error('Error applying left changes:', error)
    alert('Error applying changes')
  }
}

const applyRightChanges = async () => {
  try {
    await ApplyChanges(leftFilePath.value, rightFilePath.value, 'right')
    await compareFiles()
  } catch (error) {
    console.error('Error applying right changes:', error)
    alert('Error applying changes')
  }
}

const mergeFiles = async () => {
  try {
    const outputPath = `${leftFilePath.value}.merged`
    await MergeFiles(leftFilePath.value, rightFilePath.value, outputPath)
    alert('Files merged successfully!')
  } catch (error) {
    console.error('Error merging files:', error)
    alert('Error merging files')
  }
}
</script>

<style>
.container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  padding: 1rem;
  background-color: #f5f5f5;
}

.header {
  margin-bottom: 1rem;
}

.file-inputs {
  display: flex;
  gap: 1rem;
  margin-bottom: 1rem;
}

.file-input {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.actions {
  display: flex;
  gap: 0.5rem;
}

button {
  padding: 0.5rem 1rem;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.diff-container {
  flex: 1;
  display: flex;
  gap: 1rem;
  overflow: hidden;
}

.diff-pane {
  flex: 1;
  background-color: white;
  padding: 1rem;
  border-radius: 4px;
  overflow-y: auto;
  font-family: monospace;
}

.diff-line {
  display: flex;
  gap: 1rem;
  padding: 0.125rem 0;
}

.line-number {
  color: #666;
  min-width: 3ch;
  text-align: right;
}

.line-content {
  white-space: pre;
}

.same {
  background-color: transparent;
}

.added {
  background-color: #e6ffe6;
}

.removed {
  background-color: #ffe6e6;
}
</style>
