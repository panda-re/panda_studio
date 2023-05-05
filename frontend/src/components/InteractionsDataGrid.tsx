import { EuiBasicTable, EuiBasicTableColumn, EuiButton, EuiButtonIcon, EuiConfirmModal, EuiFlexGroup, EuiFlexItem, EuiSearchBar, EuiSearchBarOnChangeArgs, EuiSpacer, RIGHT_ALIGNMENT } from '@elastic/eui';
import { useNavigate } from 'react-router-dom';
import {InteractionProgram, Recording, useDeleteProgramById, useDeleteRecordingById, useFindAllPrograms} from "../api";
import {useQueryClient} from "@tanstack/react-query";
import {useEffect, useState} from "react";
import {useLocation} from "react-router";


function ImagesDataGrid() {
  const navigate = useNavigate();
  const {isLoading, error, data} = useFindAllPrograms();
  const location = useLocation();
  const queryClient = useQueryClient();
  const deleteFunction = useDeleteProgramById({mutation: {onSuccess: () => queryClient.invalidateQueries()}});

  const deleteProgram = ({programId}: {programId: string}) => {
    deleteFunction.mutate({programId: programId});
  }

  function deleteActionPress (item: InteractionProgram){
    deleteProgram({programId: item.id!})
    setIsConfirmVisible(false);
  }

  // Confirm Modal Fields and Methods
  const [isConfirmVisible, setIsConfirmVisible] = useState(false);
  const [itemToDelete, setItemToDelete] = useState({})

  function showConfirmModal(event: React.MouseEvent){
    setIsConfirmVisible(true);
    event.stopPropagation();
  }

  function closeConfirmModal(){
  setIsConfirmVisible(false);
  }

  function ConfirmModal(){
    return <EuiConfirmModal
      title="Are you sure you want to delete?"
      onCancel={closeConfirmModal}
      onConfirm={() => deleteActionPress(itemToDelete)}
      cancelButtonText="Cancel"
      confirmButtonText="Delete Program"
      buttonColor="danger"
      defaultFocusedButton="confirm"
    ></EuiConfirmModal>;
  }

  // Delete if program is passed back tp dashboard
  useEffect(() => {
    if(location.state) {
      deleteProgram({programId: location.state.programId});
    }
  }, []);

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
      align: RIGHT_ALIGNMENT,
      name: 'Delete',
      render: (item: InteractionProgram) => {
        return (
          <EuiButtonIcon
            onClick={(event: React.MouseEvent) => {
              setItemToDelete(item);
              showConfirmModal(event);
            }}
            iconType={"trash"}
          />
        );
      },
    },
  ]

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
    {(isConfirmVisible) ? (ConfirmModal()) : null}
  </>)
}

export default ImagesDataGrid;