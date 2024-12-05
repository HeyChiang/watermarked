// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {watermark} from '../models';
import {file} from '../models';

export function AddImageWatermark(arg1:string,arg2:string,arg3:watermark.WatermarkOptions):Promise<void>;

export function AddTextWatermark(arg1:string,arg2:watermark.WatermarkOptions):Promise<void>;

export function GetImagePreview(arg1:string):Promise<string>;

export function GetSupportedFormats():Promise<Array<string>>;

export function Greet(arg1:string):Promise<string>;

export function UploadImage(arg1:Array<number>,arg2:string):Promise<file.FileInfo>;