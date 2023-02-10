import {EuiBasicTable, EuiBasicTableColumn, EuiBasicTableProps} from '@elastic/eui';
import {useLoaderData, useLocation, useNavigate} from 'react-router';
import {Recording, useFindAllRecordings} from '../api';
import prettyBytes from 'pretty-bytes';
import {findAllRecordings} from "../api";
import type {AxiosRequestConfig} from "axios";
import ContextMenu from "./ContextMenu";
import {useMemo} from "react";
import {DefaultItemAction} from "@elastic/eui/src/components/basic_table/action_types";

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
    name: 'Actions',
    render: () => <ContextMenu/>
  }
]

const reqConfig: AxiosRequestConfig = {
  baseURL: 'http://localhost:8080/api'
}

function RecordingDataGrid() {
  const navigate = useNavigate();
  const {isLoading, error, data} = useFindAllRecordings()

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
        items={data?.data ?? []}
        rowHeader="firstName"
        columns={tableColumns}
        rowProps={getRowProps}
      />
      }
  </>)
}

export default RecordingDataGrid;