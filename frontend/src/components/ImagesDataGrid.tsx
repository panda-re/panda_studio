import { EuiBasicTable, EuiBasicTableColumn } from '@elastic/eui';
import { AxiosRequestConfig } from 'axios';
import prettyBytes from 'pretty-bytes';
import { useNavigate } from 'react-router-dom';
import { findAllImages, Image, ImageFile } from '../api';

const tableColumns: EuiBasicTableColumn<Image>[] = [
  {
    field: 'id',
    name: 'Id',
  },
  {
    field: 'name',
    name: 'File Name',
  },
  {
    field: 'files',
    name: 'Size',
    render: (value: ImageFile[]) => {
      var size = 0;
      for(var f of value){
        size+= (f.size != null) ? +f.size: 0;
      }
      return prettyBytes(size, { maximumFractionDigits: 2 });
    },
  },
]

var data: Image[] = [];

const reqConfig: AxiosRequestConfig = {
  baseURL: "http://localhost:8080/api"
}

findAllImages(reqConfig).then((value) => {
  data = value.data;
});

function ImagesDataGrid() {
  const navigate = useNavigate();
  const getRowProps = (item: Image) => {
    const id = item.id;
    return {
      'data-test-subj': `image-row-${id}`,
      onClick: () => {
        navigate('/imageDetails', {state:{item: item}})
      },
    }
  }

  return (<>
    <EuiBasicTable
      tableCaption="Images"
      items={data}
      rowHeader="firstName"
      columns={tableColumns}
      rowProps={getRowProps}
    />
  </>)
}

export default ImagesDataGrid;