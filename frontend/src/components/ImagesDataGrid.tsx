import { EuiBasicTable, EuiBasicTableColumn, EuiButtonIcon, RIGHT_ALIGNMENT } from '@elastic/eui';
import { getItemId } from '@elastic/eui/src/components/basic_table/basic_table';
import { useQueryClient } from '@tanstack/react-query';
import axios, { AxiosRequestConfig } from 'axios';
import prettyBytes from 'pretty-bytes';
import React, { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { CreateImageRequest, findAllImages, Image, ImageFile, PandaConfig, updateImage, useDeleteImageById, useFindAllImages, useUpdateImage } from '../api';

function ImagesDataGrid() {
  const navigate = useNavigate();
  const location = useLocation();
  const {isLoading, error, data} = useFindAllImages();
  const queryClient = useQueryClient();
  const deleteFunction = useDeleteImageById({mutation: {onSuccess: () => queryClient.invalidateQueries()}});
  const updateFn = useUpdateImage({mutation: {onSuccess: () => queryClient.invalidateQueries()}});

  
  const deleteImage = ({itemId}: {itemId: string}) => {
    deleteFunction.mutate({imageId: itemId});
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
        deleteImage({itemId: location.state.imageId});
      }
      window.history.replaceState({}, document.title)
    }
  }, []);

  function deleteActionPress (event: React.MouseEvent, item: Image){
    deleteImage({itemId: item.id!})
    event.stopPropagation();
  }

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
    {
      align: RIGHT_ALIGNMENT,
      name: 'Delete',
      render: (item: Image) => {
        return (
          <EuiButtonIcon
            onClick={(event: React.MouseEvent) => {deleteActionPress(event, item)}}
            iconType={"trash"}
          />
        );
      },
    },
  ]
  
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