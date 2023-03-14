import { EuiBasicTable, EuiBasicTableColumn } from '@elastic/eui';
import { useNavigate } from 'react-router-dom';
import {InteractionProgram, useDeleteProgramById, useDeleteRecordingById, useFindAllPrograms} from "../api";
import {useQueryClient} from "@tanstack/react-query";
import {useEffect} from "react";
import {useLocation} from "react-router";


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
  const location = useLocation();
  const queryClient = useQueryClient();
  const deleteFunction = useDeleteProgramById({mutation: {onSuccess: () => queryClient.invalidateQueries()}});

  const deleteProgram = ({programId}: {programId: string}) => {
    deleteFunction.mutate({programId: programId});
  }

  useEffect(() => {
    if(location.state) {
      deleteProgram({programId: location.state.programId});
    }
  }, []);

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