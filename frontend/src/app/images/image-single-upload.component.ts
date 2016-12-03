import { Component, OnInit } from '@angular/core';
import { FileUploader, FileUploaderOptions } from 'ng2-file-upload/ng2-file-upload';

import { imageUploadUrl } from 'shared/shared-data';

@Component({
	selector: 'image-single-upload',
	templateUrl: './html/image-single-upload.component.html',
	styleUrls: [ 'image-upload.component.css' ],
	providers: []
})

export class ImageSingleUploadComponent {
	private settings: FileUploaderOptions = {
		url: imageUploadUrl,
		queueLimit: 1,
		autoUpload: true,
		allowedFileType: ["jpg","png","jpeg"]
	};

	uploader: FileUploader = new FileUploader(this.settings);
	hasBaseDropZoneOver: boolean = false;

	constructor() {
	}

	ngOnInit(): void {
		// this.uploader.setOptions(this.settings);
	}

	fileOverBase(res: any): void {
		this.hasBaseDropZoneOver = res;
	}
}