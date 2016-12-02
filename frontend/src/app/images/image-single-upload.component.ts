import { Component, Input } from '@angular/core';
import { FileUploader, FileUploaderOptions } from 'ng2-file-upload/ng2-file-upload';

import { imageUploadUrl } from 'shared/shared-data';

@Component({
	selector: 'image-single-upload',
	templateUrl: './html/image-single-upload.component.html',
	styleUrls: [ 'image-upload.component.css' ],
	providers: []
})

export class ImageSingleUploadComponent {
	settings: FileUploaderOptions = {
		url: imageUploadUrl
		// url: imageUploadUrl,
		// autoUpload: true,
		// allowedMimeType: ["image"]
		// allowedFileType: ["jpg","png","jpeg"]
	};

	uploader: FileUploader = new FileUploader(this.settings);
	hasBaseDropZoneOver: boolean = false;

	@Input() criterionId: number;
	@Input() auditId: number;
	@Input() data: string;

	constructor() { }

	fileOverBase(res: any): void {
		this.hasBaseDropZoneOver = res;
	}
}
