export namespace session {
	
	export class Session {
	    apiBase: string;
	    token: string;
	
	    static createFrom(source: any = {}) {
	        return new Session(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.apiBase = source["apiBase"];
	        this.token = source["token"];
	    }
	}

}

