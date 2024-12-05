// Position type for watermark positioning
export type Position = 'center' | 'topLeft' | 'topRight' | 'bottomLeft' | 'bottomRight' | 'tiled';

// Interface for watermark options
export interface WatermarkOptions {
  text?: string;
  textSize?: number;
  textColor?: {
    R: number;
    G: number;
    B: number;
    A: number;
  };
  fontFamily?: string;
  scale?: number;
  opacity: number;
  angle: number;
  spacing: number;
  position: Position;
  margin: number;
}

// Interface for uploaded file info
export interface UploadedFile {
  name: string;
  path: string;
  size: number;
  extension: string;
}
