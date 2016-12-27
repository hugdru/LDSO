import { Component, Input, OnInit } from '@angular/core';
import { FileUploader, FileUploaderOptions } from 'ng2-file-upload/ng2-file-upload';

import { Remark } from "remark/remark";

import { imageUploadUrl } from 'shared/shared-data';

@Component({
	selector: 'image-single-upload',
	templateUrl: './html/image-single-upload.component.html',
	styleUrls: [ 'image-upload.component.css' ],
	providers: []
})

export class ImageSingleUploadComponent {
	settings: FileUploaderOptions = {
		// url: imageUploadUrl
		url: imageUploadUrl,
		// removeAfterUpload: true,
		// autoUpload: true,
		// allowedMimeType: ["image"]
		// allowedFileType: ["jpg","png","jpeg"]
	};

	uploader: FileUploader = new FileUploader(this.settings);
	hasBaseDropZoneOver: boolean = false;

	@Input() remark: Remark;

	constructor() { }

	ngOnInit(): void {
		this.uploader.onBuildItemForm = (fileItem: any, form: any) => {
			form.append('data',remark.data);
		};
	}

	upload(): void {
		for (let item of this.uploader.queue) {
			item.upload();
		}
	}

	fileOverBase(res: any): void {
		this.hasBaseDropZoneOver = res;
	}
}
