import { EuiBasicTable, EuiBasicTableColumn } from '@elastic/eui';
import { useNavigate } from 'react-router-dom';

interface InteractionProgram {
  id: string;
  name: string;
  date: Date;
};

const tableColumns: EuiBasicTableColumn<InteractionProgram>[] = [
  {
    field: 'id',
    name: 'Id',
  },
  {
    field: 'name',
    name: 'File Name',
  },
  {
    field: 'date',
    name: 'Timestamp',
  },
]

const data: InteractionProgram[] = [
  {
    id: 'IMG001',
    name: 'list-one',
    date: new Date()
  }
];

function ImagesDataGrid() {
  const navigate = useNavigate();
  const getRowProps = (item: InteractionProgram) => {
    const { id } = item;
    return {
      'data-test-subj': `image-row-${id}`,
      onClick: () => {
        navigate('/interactionDetails', {state:{item: item}})
      },
    }
  }

  return (<>
    <EuiBasicTable
      tableCaption="Interaction Programs"
      items={data}
      rowHeader="firstName"
      columns={tableColumns}
      rowProps={getRowProps}
    />
  </>)
}

export default ImagesDataGrid;