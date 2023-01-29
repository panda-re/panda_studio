import { EuiBasicTable, EuiBasicTableColumn } from '@elastic/eui';
import prettyBytes from 'pretty-bytes';
import { useNavigate } from 'react-router-dom';
import { Image } from './Interfaces';

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
    field: 'operatingSystem',
    name: 'Operating System',
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

const data: Image[] = [
  {
    id: 'AGHA68',
    name: 'wheezy.qcow2',
    operatingSystem: 'Ubuntu',
    date: new Date(),
    size: 150*1024*1024,
  }
];

function ImagesDataGrid() {
  const navigate = useNavigate();
  const getRowProps = (item: Image) => {
    const { id } = item;
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