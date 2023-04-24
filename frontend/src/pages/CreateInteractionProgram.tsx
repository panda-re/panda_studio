import {
  EuiTextArea,
  EuiButton,
  EuiButtonEmpty,
  EuiButtonIcon,
  EuiDragDropContext,
  euiDragDropCopy,
  euiDragDropReorder,
  EuiDraggable,
  EuiDroppable,
  EuiFieldSearch,
  EuiFieldText,
  EuiFlexGroup,
  EuiFlexItem,
  EuiFlyout,
  EuiFlyoutBody,
  EuiFlyoutFooter,
  EuiFlyoutHeader,
  EuiForm,
  EuiFormRow,
  EuiIcon,
  EuiModal,
  EuiModalBody,
  EuiModalFooter,
  EuiModalHeader,
  EuiModalHeaderTitle,
  EuiOverlayMask,
  EuiPageTemplate,
  EuiPanel,
  EuiSelectableOption,
  EuiSpacer,
  EuiTitle,
  htmlIdGenerator,
  EuiText
} from "@elastic/eui";
import React from "react";
import { useState } from "react";
import EntitySearchBar from "../components/EntitySearchBar";
import {CreateProgramRequest, useCreateProgram} from "../api";
import {useNavigate} from "react-router-dom";
import {useQueryClient} from "@tanstack/react-query";

function CreateInteractionProgramPage (){
  const makeId = htmlIdGenerator();
  const [instructionsValue, setInstructions] = useState('');
  const [nameValue, setName] = useState('');
  const navigate = useNavigate();

  const onInstructionsChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setInstructions(e.target.value);
  }

  const onNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setName(e.target.value);
  }

  const queryClient = useQueryClient();
  const createInteractionProgram = useCreateProgram({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries();
        navigate('/interactions')
      },
      onError: ({ response }) => alert(response?.data.error?.message),
    }
  });

  // Use to make temp data. Will replace once object storage and db implemented
  const makeList = (number: number, start = 1) =>
    Array.from({ length: number }, (v, k) => k + start).map((el) => {
      return {
        content: `Item ${el}`,
        id: makeId(),
      };
    });

  // Drag and Drop Widget constants
  const [isItemRemovable, setIsItemRemovable] = useState(false);
  const [list1, setList1] = useState(makeList(3));
  const [list2, setList2] = useState(makeList(1));
  const lists = { availableInteractions: list1, chosenInteractions: list2 };
  const actions = {
    availableInteractions: setList1,
    chosenInteractions: setList2,
  };

  // Flyout Constants
  const [isFlyoutVisible, setIsFlyoutVisible] = useState(false);
  var flyoutIndex = 0;
  const closeFlyout = () => setIsFlyoutVisible(false);
  const showFlyout = (index: number) => {
    flyoutIndex = index;
    setIsFlyoutVisible(true);
  }

  // Modal Constants
  const [isModalVisible, setIsModalVisible] = useState(false);
  var modalValue = "";
  const closeModal = () => setIsModalVisible(false);
  const showModal = () => {
    setIsModalVisible(true);
  }

  // Panel shown when no chosen interactions
  function EmptyListPanel(){
    return <EuiFlexGroup
              alignItems="center"
              justifyContent="spaceAround"
              gutterSize="none"
              style={{ height: '100%' }}>
              <EuiFlexItem grow={false}>Drop Items Here</EuiFlexItem>
            </EuiFlexGroup>
  }

  // Functions for Drag and Drop Widget
  const remove = (droppableId: string, index: number) => {
    var list = (droppableId == "availableInteractions") ? Array.from(lists.availableInteractions) : Array.from(lists.chosenInteractions);
    list.splice(index, 1);

    (droppableId == "availableInteractions") ? actions.availableInteractions(list) : actions.chosenInteractions(list);
  };

  const onDragUpdate = ({ source, destination }: {source: any, destination?: any}) => {
    const shouldRemove =
      !destination && source.droppableId == 'chosenInteractions';
    setIsItemRemovable(shouldRemove);
  };

  const onDragEnd = ({ source, destination }: {source: any, destination?: any}) => {
    if (source && destination) {
      if (source.droppableId === destination.droppableId) {
        const items = euiDragDropReorder(
          (destination.droppableId == "availableInteractions") ? lists.availableInteractions : lists.chosenInteractions,
          source.index,
          destination.index
        );
        (destination.droppableId == "availableInteractions") ? actions.availableInteractions(items) : actions.chosenInteractions(items)
      } else {
        const sourceId = source.droppableId;
        const destinationId = destination.droppableId;
        const result = euiDragDropCopy(
          (sourceId == "availableInteractions") ? lists.availableInteractions : lists.chosenInteractions,
          (destinationId == "availableInteractions") ? lists.availableInteractions : lists.chosenInteractions,
          source,
          destination,
          {
            property: 'id',
            modifier: makeId,
          }
        );
        (sourceId == "availableInteractions") ? actions.availableInteractions(result[sourceId]) : actions.chosenInteractions(result[sourceId]);
        (destinationId == "availableInteractions") ? actions.availableInteractions(result[destinationId]) : actions.chosenInteractions(result[destinationId]);
      }
    } else if (!destination && source.droppableId == 'chosenInteractions') {
      remove(source.droppableId, source.index);
    }
  };

  // Adds selected interaction from Flyout to the chosen interactions
  function addSelectedInteractionFromFlyout(){
    if(selectedInteraction != null){
      let firstHalf = list2.slice(0, flyoutIndex+1);
      let secondHalf = list2.slice(flyoutIndex+1);
      firstHalf.push({
        id: makeId(),
        content: `${selectedInteraction}`,
      });
      let list = firstHalf.concat(secondHalf);

      actions.chosenInteractions(list);
      flyoutIndex = 0;
    }
    closeFlyout();
  }

  function CreateAvailableInteractionList(){
    return list1.map(({ content, id }, idx) => (
      <EuiDraggable key={id} index={idx} draggableId={id} spacing="l">
        <EuiPanel>{content}</EuiPanel>
      </EuiDraggable>
    ))
  }

  function CreateChosenInteractionList(){
    return list2.map(({ content, id }, idx) => (
      <EuiFlexGroup gutterSize="none" direction="column">
        <EuiFlexItem>
          <EuiDraggable
            key={id}
            index={idx}
            draggableId={id}
            spacing="none"
            isRemovable={isItemRemovable}>
            <EuiPanel>
              <EuiFlexGroup gutterSize="none" alignItems="center">
                <EuiFlexItem>{content}</EuiFlexItem>
                <EuiFlexItem grow={false}>
                  {isItemRemovable ? (
                    <EuiIcon type="trash" color="danger" />
                  ) : (
                    <EuiButtonIcon
                      iconType="cross"
                      aria-label="Remove"
                      onClick={() => remove('chosenInteractions', idx)}
                    />
                  )}
                </EuiFlexItem>
              </EuiFlexGroup>
            </EuiPanel>
          </EuiDraggable>
        </EuiFlexItem>
        <EuiFlexItem grow>
          <EuiButtonEmpty
            size="xs" 
            iconType={"plusInCircle"} 
            onClick={() => {
              showFlyout(idx);
            }}>
          </EuiButtonEmpty>
        </EuiFlexItem>
      </EuiFlexGroup>
    ))
  }

  // Creates the Drag and Drop Widget
  function CreateDroppableWidget() {
    return (<>
      <EuiDragDropContext onDragEnd={onDragEnd} onDragUpdate={onDragUpdate}>
        <EuiFlexGroup>
          <EuiFlexItem>
            <EuiFlexItem >
              <EuiDroppable
                droppableId="availableInteractions"
                cloneDraggables={true}
                grow
                withPanel>
                {CreateAvailableInteractionList()}
              </EuiDroppable>
            </EuiFlexItem>
          </EuiFlexItem>
          <EuiFlexItem>
            <EuiDroppable droppableId="chosenInteractions" withPanel grow>
              <EuiFlexGroup  gutterSize="none" direction="column">
                <EuiFlexItem>
                  {list2.length ? (CreateChosenInteractionList()) : (<EmptyListPanel></EmptyListPanel>)}
                </EuiFlexItem>
              </EuiFlexGroup>
            </EuiDroppable>
          </EuiFlexItem>
        </EuiFlexGroup>
      </EuiDragDropContext>
    </>);
  };

  // Generate selectable options for add Interaction flyout search component
  let interactionEntities: EuiSelectableOption[] = [];
  list1.map((i) =>
    interactionEntities.push({label: i.content,
    data: i})
  );
  const [selectedInteraction, setSelectedInteraction] = React.useState<EuiSelectableOption | undefined>(undefined);
  function returnSelectedInteraction(interaction: EuiSelectableOption){
   setSelectedInteraction(interaction);
  }

  function CreateFlyout(){
    return <EuiFlyout onClose={closeFlyout}>
            <EuiFlyoutHeader>
              <EuiTitle>
                <h2>Add Interaction </h2>
              </EuiTitle>
            </EuiFlyoutHeader>
            <EuiFlyoutBody>
              <EuiPanel>
                <EntitySearchBar name="Interaction" entities={interactionEntities} returnSelectedOption={(returnSelectedInteraction)}></EntitySearchBar>
              </EuiPanel>
            </EuiFlyoutBody>
            <EuiFlyoutFooter>
              <EuiFlexGroup justifyContent="spaceBetween">
                <EuiFlexItem grow={false}>
                  <EuiButtonEmpty onClick={closeFlyout}>Cancel</EuiButtonEmpty>
                </EuiFlexItem>
                <EuiFlexItem grow={false}>
                  <EuiButton onClick={addSelectedInteractionFromFlyout} fill>Save</EuiButton>
                </EuiFlexItem>
              </EuiFlexGroup>
            </EuiFlyoutFooter>
          </EuiFlyout>
  }

  function CreateModal(){
    return <EuiOverlayMask>
              <EuiModal onClose={closeModal}>
                <EuiModalHeader>
                  <EuiModalHeaderTitle>Add New Interaction</EuiModalHeaderTitle>
                </EuiModalHeader>
                <EuiModalBody>
                    <EuiFieldText 
                      placeholder="Enter New Interaction"  
                      name="New Interaction" 
                      onChange={(e) => {
                        modalValue = e.target.value;
                      }}/>
                </EuiModalBody>
                <EuiModalFooter>
                  <EuiButton onClick={closeModal} fill>Close</EuiButton>
                  <EuiButton 
                    onClick={() => {
                      let list = list1;
                      list.push({
                        id: makeId(),
                        content: modalValue,
                      })
                      setList1(list);
                      closeModal();
                    }} 
                    fill>
                      Submit</EuiButton>
                </EuiModalFooter>
              </EuiModal>
            </EuiOverlayMask>
  }

  return (<>
    <EuiPageTemplate.Header pageTitle='Create Interaction Program' />
    <EuiPageTemplate.Section>
      <EuiFlexGroup justifyContent={"spaceEvenly"}>
        <EuiFlexItem grow={3}>
          <EuiFieldText
            placeholder={"Program Name"}
            value={nameValue}
            onChange={e => onNameChange(e)}>
          </EuiFieldText>
        </EuiFlexItem>

        <EuiFlexItem grow={6}>
          <div>
            <EuiButton onClick={() => {
              if (instructionsValue == "") {
                alert('No instructions?')
                return
              }
              const createProgramRequest: CreateProgramRequest = {
                name: nameValue,
                instructions: instructionsValue
              }
              createInteractionProgram.mutate({data: createProgramRequest})
            }}>Create Interaction Program</EuiButton>
          </div>
        </EuiFlexItem>
      </EuiFlexGroup>
      <EuiSpacer size="m"></EuiSpacer>
      {(isFlyoutVisible) ? (CreateFlyout()) : null}
      {(isModalVisible) ? (CreateModal()) : null}
      <EuiTextArea
        fullWidth={true}
        value={instructionsValue}
        onChange={e => onInstructionsChange(e)}>
      </EuiTextArea>
    </EuiPageTemplate.Section>
  </>);
}

export default CreateInteractionProgramPage;