import { EuiSelectable, EuiSelectableOption } from '@elastic/eui';
import React, { useState, Fragment, MutableRefObject } from 'react';

// name: used in placeholder of search bar
// entities: list of "stringified" entities to let user select from
// returnSelected: callback function (hook) defined by parent component to retrieve the selected item from the child component
// ------------------------------------------------------------
//    create a callback like this in the parent's component construction method: 
// const [selectedRecording, setSelectedRecording] = React.useState<String | undefined>("NONE SELECTED MESSAGE");
// function returnSelected(message: String){
//    setSelectedRecording(message);
// }
interface EntitySearchBarProps{
    name: String;
    entities: EuiSelectableOption[]
    returnSelectedOption: Function,
}

function EntitySearchBar(props: EntitySearchBarProps) {
    const [entities, setOptions] = useState(props.entities);
    return (<>
        <Fragment>
            <EuiSelectable
                searchable
                singleSelection={true}
                searchProps={{
                    'data-test-subj': 'selectableSearchHere',
                    'placeholder': 'Filter ' + props.name,
                }}
                options={entities}
                onChange={(newOptions) => {
                    setOptions(newOptions);
                    const hasSelected = newOptions.find(selected => selected.checked === "on");
                    props.returnSelectedOption(hasSelected !== undefined ? hasSelected : "No Recording Selected");
                }}>
                {(list, search) => (
                    <Fragment>
                        {search}
                        {list}
                    </Fragment>
                )}
            </EuiSelectable>
        </Fragment>
    </>)
}

export default EntitySearchBar;