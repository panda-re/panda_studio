import React, { useState } from 'react';
import {
  EuiButton,
  EuiContextMenuPanel,
  EuiContextMenuItem,
  EuiPopover,
  EuiCopy,
  useGeneratedHtmlId, EuiButtonEmpty,
} from '@elastic/eui';
import {Recording, useDeleteRecordingById} from "../api";
import {useQueryClient} from "@tanstack/react-query";

export default ({recordingId, deleteCallback}: {recordingId: string, deleteCallback: any}) => {
  const [isPopoverOpen, setPopover] = useState(false);
  const smallContextMenuPopoverId = useGeneratedHtmlId({
    prefix: 'smallContextMenuPopover',
  });

  const onButtonClick: React.MouseEventHandler = (event ) => {
    setPopover(!isPopoverOpen);
    event.stopPropagation();
  };

  const closePopover = () => {
    setPopover(false);
  };

  const deleteItem: React.MouseEventHandler = (event) => {
    deleteCallback({recordingId: recordingId});
    closePopover();
    event.stopPropagation();
  };

  const items = [
    <EuiContextMenuItem key="delete" icon="trash" onClick={deleteItem}>
      Delete
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