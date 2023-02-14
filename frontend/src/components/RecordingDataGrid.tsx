import {EuiBasicTable, EuiBasicTableColumn, EuiBasicTableProps} from '@elastic/eui';
import {useLoaderData, useLocation, useNavigate} from 'react-router';
import {Recording, useFindAllRecordings} from '../api';
import prettyBytes from 'pretty-bytes';
import ContextMenu from "./ContextMenu";

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
]

function RecordingDataGrid() {
  const navigate = useNavigate();
  const {isLoading, error, data} = useFindAllRecordings();

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
        render: (item: Recording) => <ContextMenu recordingId={item.id!} />
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