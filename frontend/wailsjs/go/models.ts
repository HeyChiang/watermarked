export namespace color {
	
	export class RGBA {
	    R: number;
	    G: number;
	    B: number;
	    A: number;
	
	    static createFrom(source: any = {}) {
	        return new RGBA(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.R = source["R"];
	        this.G = source["G"];
	        this.B = source["B"];
	        this.A = source["A"];
	    }
	}

}

export namespace file {
	
	export class FileInfo {
	    name: string;
	    path: string;
	    size: number;
	    extension: string;
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.size = source["size"];
	        this.extension = source["extension"];
	    }
	}

}

export namespace watermark {
	
	export class WatermarkOptions {
	    text: string;
	    textSize: number;
	    textColor: color.RGBA;
	    fontFamily: string;
	    scale: number;
	    opacity: number;
	    angle: number;
	    spacing: number;
	    position: string;
	    margin: number;
	
	    static createFrom(source: any = {}) {
	        return new WatermarkOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.text = source["text"];
	        this.textSize = source["textSize"];
	        this.textColor = this.convertValues(source["textColor"], color.RGBA);
	        this.fontFamily = source["fontFamily"];
	        this.scale = source["scale"];
	        this.opacity = source["opacity"];
	        this.angle = source["angle"];
	        this.spacing = source["spacing"];
	        this.position = source["position"];
	        this.margin = source["margin"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

