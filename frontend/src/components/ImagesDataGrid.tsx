import { EuiBasicTable, EuiBasicTableColumn } from '@elastic/eui';

import prettyBytes from 'pretty-bytes';

// image ID, name, OS, Timestamp, Size, view specs
interface Image {
  id: string;
  name: string;
  date: Date;
  operatingSystem: string;
  size: number;
};

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

const getRowProps = (item: Image) => {
  const { id } = item;
  return {
    'data-test-subj': `image-row-${id}`,
    onClick: () => {
      alert(`View Image with id ${id}`)
    },
  }
};

function ImagesDataGrid() {

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