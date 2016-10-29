import { Property } from './property';
import { PropertyService } from './service/property.service';

@Component({
	moduleId: module.id,
	selector: 'properties-info',
	templateUrl: './html/properties-info.component.html',
	styleUrls: [ './css/properties-info.component.css' ],
	providers: [ PropertyService ]
})

export class PropertiesInfoComponent implements OnInit {

}
