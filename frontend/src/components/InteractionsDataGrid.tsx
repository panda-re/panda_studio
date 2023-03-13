import { EuiBasicTable, EuiBasicTableColumn } from '@elastic/eui';
import { useNavigate } from 'react-router-dom';
import {InteractionProgram, useFindAllPrograms} from "../api";


const tableColumns: EuiBasicTableColumn<InteractionProgram>[] = [
  {
    field: 'id',
    name: 'Id',
  },
  {
    field: 'name',
    name: 'File Name',
  },
]


function ImagesDataGrid() {
  const navigate = useNavigate();
  const {isLoading, error, data} = useFindAllPrograms();

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
    {isLoading && <div>Loading...</div> ||
    <EuiBasicTable
      tableCaption="Interaction Programs"
      items={data ?? []}
      rowHeader="firstName"
      columns={tableColumns}
      rowProps={getRowProps}
    />
    }
  </>)
}

export default ImagesDataGrid;