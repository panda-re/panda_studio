import { EuiBasicTable, EuiBasicTableColumn } from '@elastic/eui';
import { useQueryClient } from '@tanstack/react-query';
import axios, { AxiosRequestConfig } from 'axios';
import prettyBytes from 'pretty-bytes';
import { useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { findAllImages, Image, ImageFile, useDeleteImageById, useFindAllImages } from '../api';

function ImagesDataGrid() {
  const location = useLocation();
  const queryClient = useQueryClient();
  const deleteFunction = useDeleteImageById({mutation: {onSuccess: () => queryClient.invalidateQueries()}});

  const {isLoading, error, data} = useFindAllImages();

  const deleteImage = ({imageId}: {imageId: string}) => {
    deleteFunction.mutate({imageId: imageId});
  }

  useEffect(() => {
    if(location.state) {
      deleteImage({imageId: location.state.recordingId});
    }
  }, []);

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