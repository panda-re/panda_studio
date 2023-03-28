import {EuiBasicTable, EuiBasicTableColumn, EuiBasicTableProps} from '@elastic/eui';
import {useLoaderData, useLocation, useNavigate} from 'react-router';
import {Recording, useDeleteRecordingById, useFindAllRecordings} from '../api';
import prettyBytes from 'pretty-bytes';
import ContextMenu from "./ContextMenu";
import {useEffect, useState} from "react";
import {useQueryClient} from "@tanstack/react-query";

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
  /*
  {
    field: 'size',
    name: 'Size',
    render: (value: number) => prettyBytes(value, {maximumFractionDigits: 2}),
  },
  */
  {
    field: 'date',
    name: 'Timestamp',
  },
]

function RecordingDataGrid() {
  const navigate = useNavigate();
  const location = useLocation();
  const {isLoading, error, data} = useFindAllRecordings();
  const queryClient = useQueryClient();
  const deleteFunction = useDeleteRecordingById({mutation: {onSuccess: () => queryClient.invalidateQueries()}});

  const deleteRecording = ({recordingId}: {recordingId: string}) => {
    deleteFunction.mutate({recordingId: recordingId});
  }

  useEffect(() => {
    if(location.state) {
      deleteRecording({recordingId: location.state.recordingId});
    }
  }, []);

  const getRowProps: EuiBasicTableProps<Recording>['rowProps'] = (item) => {
    const {id} = item;
    return {
      'data-test-subj': `recording-row-${id}`,
      onClick: () => {
        navigate('/recordingDetails', {state: {item: item}});
      },
    }
  };

  const columnsWithActions = [
    ...tableColumns,
    {
      name: 'Actions',
        render: (item: Recording) => <ContextMenu recordingId={item.id!} deleteCallback={deleteRecording} />
    },
  ]

  return (<>
    {isLoading && <div>Loading...</div> ||
      <EuiBasicTable
        tableCaption="Recordings"
        items={data ?? []}
        rowHeader="firstName"
        columns={columnsWithActions}
        rowProps={getRowProps}
      />
      }
  </>)
}

export default RecordingDataGrid;