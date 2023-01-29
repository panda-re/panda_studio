import { EuiSelectable, EuiSelectableOption } from '@elastic/eui';
import React, { useState, Fragment, MutableRefObject } from 'react';

// name: used in placeholder of search bar
// entities: list of "stringified" entities to let user select from
// returnSelected: callback function (hook) defined by parent component to retrieve the selected item from the child component
// ------------------------------------------------------------
// let imageEntities: EuiSelectableOption[] = [];
//   data.map((r) =>
//     imageEntities.push({label: `${r.id} - ${r.name} - ${r.operatingSystem} - ${r.date.toLocaleDateString()} - ${prettyBytes(r.size, { maximumFractionDigits: 2 })}`,
//   data: r})
//   );

//   const [selectedImage, setSelectedImage] = React.useState<EuiSelectableOption | undefined>(undefined);
//   function returnSelectedImage(message: EuiSelectableOption){
//    setSelectedImage(message);
//   }
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