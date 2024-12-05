import { ApplyChanges, MergeFiles, CompareFiles as GoCompareFiles } from '../wailsjs/go/main/App'

export type DiffLineType = 'same' | 'added' | 'removed'

export interface DiffLine {
  lineNumber: number
  content: string
  type: DiffLineType
}

export interface DiffResult {
  leftLines: DiffLine[]
  rightLines: DiffLine[]
}

export async function CompareFiles(leftPath: string, rightPath: string): Promise<DiffResult> {
  const result = await GoCompareFiles(leftPath, rightPath)
  return result as DiffResult
}
