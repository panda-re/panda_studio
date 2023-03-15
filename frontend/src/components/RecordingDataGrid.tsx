import {EuiBasicTable, EuiBasicTableColumn, EuiBasicTableProps, EuiButtonIcon, RIGHT_ALIGNMENT} from '@elastic/eui';
import {useLoaderData, useLocation, useNavigate} from 'react-router';
import {Recording, useDeleteRecordingById, useFindAllRecordings} from '../api';
import prettyBytes from 'pretty-bytes';
import {useEffect, useState} from "react";
import {useQueryClient} from "@tanstack/react-query";

function RecordingDataGrid() {
  const navigate = useNavigate();
  const location = useLocation();
  const {isLoading, error, data} = useFindAllRecordings();
  const queryClient = useQueryClient();
  const deleteFunction = useDeleteRecordingById({mutation: {onSuccess: () => queryClient.invalidateQueries()}});

  const deleteRecording = ({itemId}: {itemId: string}) => {
    deleteFunction.mutate({recordingId: itemId});
  }

  useEffect(() => {
    if(location.state) {
      deleteRecording({itemId: location.state.recordingId});
    }
  }, []);

  function deleteActionPress (event: React.MouseEvent, item: Recording){
    deleteRecording({itemId: item.id!})
    event.stopPropagation();
  }

  const tableColumns: EuiBasicTableColumn<Recording>[] = [
    {
      field: 'id',
      name: 'Id',
    },
    {
      field: 'name',
      name: 'File Name',
    },
    {
      field: 'recordingImage',
      name: 'Image Name',
    },
    {
      field: 'size',
      name: 'Size',
      render: (value: number) => prettyBytes(value, {maximumFractionDigits: 2}),
    },
    {
      field: 'date',
      name: 'Timestamp',
    },
    {
      align: RIGHT_ALIGNMENT,
      name: 'Delete',
      render: (item: Recording) => {
        return (
          <EuiButtonIcon
            onClick={(event: React.MouseEvent) => {deleteActionPress(event, item)}}
            iconType={"trash"}
          />
        );
      },
    },
  ]

  const getRowProps: EuiBasicTableProps<Recording>['rowProps'] = (item) => {
    const {id} = item;
    return {
      'data-test-subj': `recording-row-${id}`,
      onClick: () => {
        navigate('/recordingDetails', {state: {item: item}});
      },
    }
  };

  return (<>
    {isLoading && <div>Loading...</div> ||
      <EuiBasicTable
        tableCaption="Recordings"
        items={data ?? []}
        rowHeader="firstName"
        columns={tableColumns}
        rowProps={getRowProps}
      />
      }
  </>)
}

export default RecordingDataGrid;