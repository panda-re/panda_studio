import { EuiBasicTable, EuiBasicTableColumn } from '@elastic/eui';
import { useQueryClient } from '@tanstack/react-query';
import axios, { AxiosRequestConfig } from 'axios';
import prettyBytes from 'pretty-bytes';
import { useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { CreateImageRequest, findAllImages, Image, ImageFile, PandaConfig, updateImage, useDeleteImageById, useFindAllImages, useUpdateImage } from '../api';

function ImagesDataGrid() {
  const location = useLocation();
  const queryClient = useQueryClient();
  const deleteFunction = useDeleteImageById({mutation: {onSuccess: () => queryClient.invalidateQueries()}});
  const updateFn = useUpdateImage({mutation: {onSuccess: () => queryClient.invalidateQueries()}});

  const {isLoading, error, data} = useFindAllImages();

  const deleteImage = ({imageId}: {imageId: string}) => {
    deleteFunction.mutate({imageId: imageId});
  }

  const updateImage = ({image}: {image: Image}) => {
    if(image.id == null){
      return;
    }
    const conf: PandaConfig = {
      key: image.config,
    }
    const req: CreateImageRequest = {
      name: image.name,
      description: image.description,
      config: conf,
    };
    updateFn.mutate({data: req, imageId: image.id});
  }

  useEffect(() => {
    if(location.state) {
      if(location.state.image){
        updateImage({image: location.state.image});
      }
      else{
        deleteImage({imageId: location.state.imageId});
      }
      window.history.replaceState({}, document.title)
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