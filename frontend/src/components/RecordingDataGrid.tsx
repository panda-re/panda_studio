import { EuiBasicTable, EuiBasicTableColumn } from '@elastic/eui';
import { useNavigate } from 'react-router';

import prettyBytes from 'pretty-bytes';

// Recording ID, name, OS, Timestamp, Size, view specs
interface Recording {
  id: string;
  name: string;
  date: Date;
  imageName: string;
  size: number;
};

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
    field: 'imageName',
    name: 'Image Name',
  },
  {
    field: 'size',
    name: 'Size',
    render: (value: number) => prettyBytes(value, { maximumFractionDigits: 2 }),
  },
  {
    field: 'date',
    name: 'Timestamp',
  },
]

const data: Recording[] = [
  {
    id: 'record_1',
    name: 'test_recording',
    imageName: 'wheezy.qcow2',
    date: new Date(),
    size: 150*1024*1024,
  },
  {
    id: 'record_2',
    name: 'test_recording2',
    imageName: 'wheezy.qcow2',
    date: new Date(),
    size: 150*1024*1024,
  }
];

function RecordingDataGrid() {
  const navigate = useNavigate();

  const getRowProps = (item: Recording) => {
    const { id } = item;
    return {
      'data-test-subj': `recording-row-${id}`,
      onClick: () => {
        navigate('/recordingDetails', {state:{item: item}});
      },
    }
  };

  return (<>
    <EuiBasicTable
      tableCaption="Recordings"
      items={data}
      rowHeader="firstName"
      columns={tableColumns}
      rowProps={getRowProps}
    />
  </>)
}

export default RecordingDataGrid;