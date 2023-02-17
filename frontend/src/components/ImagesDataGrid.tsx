import { EuiBasicTable, EuiBasicTableColumn } from '@elastic/eui';
import axios, { AxiosRequestConfig } from 'axios';
import prettyBytes from 'pretty-bytes';
import { useNavigate } from 'react-router-dom';
import { findAllImages, Image, ImageFile, useFindAllImages } from '../api';

function ImagesDataGrid() {

  const {isLoading, error, data} = useFindAllImages();

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
    {isLoading && <div>Loading...</div> ||
    <EuiBasicTable
      tableCaption="Images"
      items={data ?? []}
      rowHeader="firstName"
      columns={tableColumns}
      rowProps={getRowProps}
    />
  }
  </>)
}

export default ImagesDataGrid;