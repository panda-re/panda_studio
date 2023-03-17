import { EuiBasicTable, EuiBasicTableColumn, EuiButton, EuiFlexGroup, EuiFlexItem, EuiSearchBar, EuiSearchBarOnChangeArgs, EuiSpacer } from '@elastic/eui';
import { useNavigate } from 'react-router-dom';
import {InteractionProgram, Recording, useDeleteProgramById, useDeleteRecordingById, useFindAllPrograms} from "../api";
import {useQueryClient} from "@tanstack/react-query";
import {useEffect, useState} from "react";
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

  const initialQuery = EuiSearchBar.Query.MATCH_ALL;

  const [query, setQuery] = useState(initialQuery);

  const onChange = (args: EuiSearchBarOnChangeArgs) => {
    setQuery(args.query ?? initialQuery);
  };

  const queriedItems = EuiSearchBar.Query.execute(query, data ?? []);

  return (<>
  <EuiFlexGroup justifyContent='spaceBetween'>
      <EuiFlexItem grow={false} style={{ minWidth: 300 }}>
        <EuiSearchBar 
          box={{
            incremental: true,
          }}
          defaultQuery={initialQuery}
          onChange={onChange}/>
      </EuiFlexItem>
      <EuiFlexItem grow={false}>
          <EuiButton iconType={'plusInCircle'} onClick={() => navigate('/createInteractionProgram')}>Create New Interaction</EuiButton>
        </EuiFlexItem>
    </EuiFlexGroup>
    <EuiSpacer></EuiSpacer>
    {isLoading && <div>Loading...</div> ||
    <EuiBasicTable
      tableCaption="Interaction Programs"
      items={queriedItems ?? []}
      rowHeader="firstName"
      columns={tableColumns}
      rowProps={getRowProps}
    />
    }
  </>)
}

export default ImagesDataGrid;