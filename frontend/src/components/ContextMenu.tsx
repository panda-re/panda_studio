import React, { useState } from 'react';
import {
  EuiButton,
  EuiContextMenuPanel,
  EuiContextMenuItem,
  EuiPopover,
  EuiCopy,
  useGeneratedHtmlId, EuiButtonEmpty,
} from '@elastic/eui';
import {deleteRecordingById, useDeleteRecordingById, useFindAllRecordings} from "../api";
import axios from "axios";

export default ({recordingId}: {recordingId: string}) => {
  const [isPopoverOpen, setPopover] = useState(false);
  const smallContextMenuPopoverId = useGeneratedHtmlId({
    prefix: 'smallContextMenuPopover',
  });

  const deleteRecording = useDeleteRecordingById();

  const onButtonClick: React.MouseEventHandler = (event ) => {
    setPopover(!isPopoverOpen);
    event.stopPropagation();
  };

  const closePopover = () => {
    setPopover(false);
  };

  const deleteItem: React.MouseEventHandler = (event) => {
    console.log(recordingId);
    deleteRecording.mutate({recordingId});
    closePopover();
    event.stopPropagation();
  };

  const items = [
    <EuiContextMenuItem key="delete" icon="trash" onClick={deleteItem}>
      Delete
    </EuiContextMenuItem>,
    <EuiContextMenuItem key="edit" icon="pencil">
      Edit
    </EuiContextMenuItem>,
    <EuiContextMenuItem key="share" icon="copy">
      Share
    </EuiContextMenuItem>,
  ];

  const button = (
    <EuiButtonEmpty color={"text"} iconType="boxesVertical" iconSide="right" onClick={onButtonClick}>
    </EuiButtonEmpty>
  );

  return (
    <EuiPopover
      id={smallContextMenuPopoverId}
      button={button}
      isOpen={isPopoverOpen}
      closePopover={closePopover}
      panelPaddingSize="none"
      anchorPosition="downLeft"
    >
      <EuiContextMenuPanel size="s" items={items} />
    </EuiPopover>
  );
}