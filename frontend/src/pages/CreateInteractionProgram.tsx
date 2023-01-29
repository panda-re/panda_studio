import { EuiButtonIcon, EuiDragDropContext, EuiDraggable, EuiDroppable, EuiFlexGroup, EuiFlexItem, EuiIcon, EuiPageTemplate, EuiPanel, htmlIdGenerator } from "@elastic/eui";
import { useState } from "react";
import {
  euiDragDropCopy,
  euiDragDropReorder,
} from '@elastic/eui';

const makeId = htmlIdGenerator();

const makeList = (number: number, start = 1) =>
  Array.from({ length: number }, (v, k) => k + start).map((el) => {
    return {
      content: `Item ${el}`,
      id: makeId(),
    };
  });

const data = [
  {
    content: "Test Item",
    id: makeId(),
  }
]

function CreateDroppableWidget() {
  const [isItemRemovable, setIsItemRemovable] = useState(false);
  const [list1, setList1] = useState(makeList(3));
  const [list2, setList2] = useState(data);
  const lists = { availableInteractions: list1, chosenInteractions: list2 };
  const actions = {
    availableInteractions: setList1,
    chosenInteractions: setList2,
  };

  const remove = (droppableId: string, index: number) => {
    var list = (droppableId == "availableInteractions") ? Array.from(lists.availableInteractions) : Array.from(lists.chosenInteractions);
    list.splice(index, 1);

    (droppableId == "availableInteractions") ? actions.availableInteractions(list) : actions.chosenInteractions(list);
  };

  const onDragUpdate = ({ source, destination }: {source: any, destination: any}) => {
    const shouldRemove =
      !destination && source.droppableId == 'chosenInteractions';
    setIsItemRemovable(shouldRemove);
  };

  const onDragEnd = ({ source, destination }: {source: any, destination: any}) => {
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

  return (
    <EuiDragDropContext onDragEnd={onDragEnd} onDragUpdate={onDragUpdate}>
      <EuiFlexGroup>
        <EuiFlexItem style={{ width: '50%' }}>
          <EuiDroppable
            droppableId="availableInteractions"
            cloneDraggables={true}
            spacing="l"
            grow>
            {list1.map(({ content, id }, idx) => (
              <EuiDraggable key={id} index={idx} draggableId={id} spacing="l">
                <EuiPanel>{content}</EuiPanel>
              </EuiDraggable>
            ))}
          </EuiDroppable>
        </EuiFlexItem>
        <EuiFlexItem style={{ width: '50%' }}>
          <EuiDroppable droppableId="chosenInteractions" withPanel grow>
            {list2.length ? (
              list2.map(({ content, id }, idx) => (
                <EuiDraggable
                  key={id}
                  index={idx}
                  draggableId={id}
                  spacing="l"
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
              ))
            ) : (
              <EuiFlexGroup
                alignItems="center"
                justifyContent="spaceAround"
                gutterSize="none"
                style={{ height: '100%' }}>
                <EuiFlexItem grow={false}>Drop Items Here</EuiFlexItem>
              </EuiFlexGroup>
            )}
          </EuiDroppable>
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiDragDropContext>
  );
};

function CreateInteractionProgramPage (){
  return (<>
    <EuiPageTemplate.Header pageTitle='Create Interaction Program' />
    <EuiPageTemplate.Section>
      <CreateDroppableWidget></CreateDroppableWidget>
    </EuiPageTemplate.Section>

  </>);
}

export default CreateInteractionProgramPage;