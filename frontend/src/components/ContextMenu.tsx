import React, { useState } from 'react';
import {
  EuiButton,
  EuiContextMenuPanel,
  EuiContextMenuItem,
  EuiPopover,
  EuiCopy,
  useGeneratedHtmlId, EuiButtonEmpty,
} from '@elastic/eui';

export default () => {
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

  const items = [
    <EuiContextMenuItem key="delete" icon="trash" onClick={closePopover}>
      Delete
    </EuiContextMenuItem>,
    <EuiContextMenuItem key="edit" icon="pencil" onClick={closePopover}>
      Edit
    </EuiContextMenuItem>,
    <EuiContextMenuItem key="share" icon="copy" onClick={closePopover}>
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