import {EuiBasicTable, EuiBasicTableColumn, EuiBasicTableProps, EuiButton, EuiButtonIcon, EuiConfirmModal, EuiFlexGroup, EuiFlexItem, EuiSearchBar, EuiSearchBarOnChangeArgs, EuiSpacer, formatDate, RIGHT_ALIGNMENT} from '@elastic/eui';
import {useLoaderData, useLocation, useNavigate} from 'react-router';
import {Recording, useDeleteRecordingById, useFindAllRecordings, useFindImageById} from '../api';
import prettyBytes from 'pretty-bytes';
import {useEffect, useState} from "react";
import {useQueryClient} from "@tanstack/react-query";
import moment from 'moment';

function RecordingDataGrid() {
  const navigate = useNavigate();
  const location = useLocation();
  const {isLoading, error, data} = useFindAllRecordings();
  const queryClient = useQueryClient();
  const deleteFunction = useDeleteRecordingById({mutation: {onSuccess: () => queryClient.invalidateQueries()}});

  const deleteRecording = ({itemId}: {itemId: string}) => {
    deleteFunction.mutate({recordingId: itemId});
  }

  useEffect(() => {
    if(location.state) {
      deleteRecording({itemId: location.state.recordingId});
    }
  }, []);

  function deleteActionPress (item: Recording){
    deleteRecording({itemId: item.id!})
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
      confirmButtonText="Delete Recording"
      buttonColor="danger"
      defaultFocusedButton="confirm"
    ></EuiConfirmModal>;
  }

  const tableColumns: EuiBasicTableColumn<Recording>[] = [
    {
      field: 'id',
      name: 'Id',
    },
    {
      field: 'name',
      name: 'Recording Name',
    },
    {
      field: 'image_id',
      name: 'Image Id',
    },
    {
      field: 'program_id',
      name: 'Program Id',
    },
    {
      field: 'date',
      name: 'Date',
      render: (value: string) => {
        return formatDate(moment(value.slice(0, 19)), 'dateTime')
      }
    },
    {
      align: RIGHT_ALIGNMENT,
      name: 'Delete',
      render: (item: Recording) => {
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

  const getRowProps: EuiBasicTableProps<Recording>['rowProps'] = (item) => {
    const {id} = item;
    return {
      'data-test-subj': `recording-row-${id}`,
      onClick: () => {
        navigate('/recordingDetails', {state: {item: item}});
      },
    }
  };

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
          <EuiButton iconType={'plusInCircle'} onClick={() => navigate('/createRecording')}>Create New Recording</EuiButton>
        </EuiFlexItem>
    </EuiFlexGroup>
    <EuiSpacer></EuiSpacer>
    {isLoading && <div>Loading...</div> ||
      <EuiBasicTable
        tableCaption="Recordings"
        items={queriedItems ?? []}
        rowHeader="firstName"
        columns={tableColumns}
        rowProps={getRowProps}
      />
      }
      {(isConfirmVisible) ? (ConfirmModal()) : null}
  </>)
}

export default RecordingDataGrid;