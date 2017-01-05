import { Component, Input, OnInit } from '@angular/core';
import { FileUploader, FileUploaderOptions }
		from 'ng2-file-upload/ng2-file-upload';

import { Remark } from "remark/remark";

import { imageUploadUrl } from 'shared/shared-data';

@Component({
	selector: 'image-single-upload',
	templateUrl: './html/image-single-upload.component.html',
	styleUrls: [ 'image-upload.component.css' ],
	providers: []
})

export class ImageSingleUploadComponent {
	imageUrl; string;
	settings: FileUploaderOptions;
	// settings: FileUploaderOptions = {
	// 	// url: imageUploadUrl
	// 	url: imageUploadUrl.replace(/#/g, this.remark.idAudit.toString())
	// 			.replace(/!/g, this.remark.idCriterion.toString())
	// 	// removeAfterUpload: true,
	// 	// autoUpload: true,
	// 	// allowedMimeType: ["image"]
	// 	// allowedFileType: ["jpg","png","jpeg"]
	// };

	// uploader: FileUploader = new FileUploader(this.settings);
	uploader: FileUploader;

	@Input() remark: Remark;

	constructor() { }
	
	ngOnInit(): void {
		this.settings = {
			url: imageUploadUrl.replace(/#/g, this.remark.idAudit.toString())
					.replace(/!/g, this.remark.idCriterion.toString())
		};
		this.uploader = new FileUploader(this.settings);
		this.uploader.onBuildItemForm = (item, form) => {
			form.append("observation", this.remark.data);
		};		
	}
}
