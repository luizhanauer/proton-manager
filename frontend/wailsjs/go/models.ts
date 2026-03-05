export namespace main {
	
	export class ProtonRelease {
	    version: string;
	    url: string;
	    size: number;
	    date: string;
	    major: string;
	    installed: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ProtonRelease(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.url = source["url"];
	        this.size = source["size"];
	        this.date = source["date"];
	        this.major = source["major"];
	        this.installed = source["installed"];
	    }
	}

}

