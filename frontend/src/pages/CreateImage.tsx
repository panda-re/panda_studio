import {EuiButton, EuiPageTemplate, EuiText, EuiFlexGroup, EuiFlexItem, EuiFieldText, EuiSuperSelect, EuiHealth, EuiSelect, EuiFilePicker, EuiFieldNumber} from '@elastic/eui';
import React, {ChangeEvent, SetStateAction, useState} from 'react';
import { useLocation } from 'react-router';

function CreateImagePage() {
    const location = useLocation()

    //hook for getting size option for image size
    const sizeOptions = [
        { value: 'kb', text: 'KB' },
        { value: 'mb', text: 'MB' },
        { value: 'gb', text: 'GB' },
    ];
    const [value, setValue] = useState(sizeOptions[2].value);
    function onSizeChange(newSizeValue: ChangeEvent<HTMLSelectElement>){
        setValue(newSizeValue.target.value)
    }

    return (<>
        <EuiPageTemplate.Header pageTitle={"Derive Image from " + location.state.item.name}/>

        <EuiPageTemplate.Section>
            <EuiFlexGroup>
                <EuiFlexItem grow={2}>
                    <EuiText>New Image Name: </EuiText>
                </EuiFlexItem>
                <EuiFlexItem grow={8}>
                    <EuiFieldText
                        placeholder="eg, image1"
                    />
                </EuiFlexItem>
            </EuiFlexGroup>
        </EuiPageTemplate.Section>
        
        <EuiPageTemplate.Section>
            <EuiFlexGroup>
                <EuiFlexItem grow={2}>
                    <EuiText>Size: </EuiText>
                </EuiFlexItem>

                <EuiFlexItem grow={8}>
                    <EuiFlexGroup>
                        <EuiFlexItem >
                            <EuiFieldText/>
                        </EuiFlexItem>
                        <EuiFlexItem>
                            <EuiSelect
                                options={sizeOptions}
                                value={value}
                                onChange={value => onSizeChange(value) }
                            />
                        </EuiFlexItem>
                    </EuiFlexGroup>
                </EuiFlexItem>     
            </EuiFlexGroup>
        </EuiPageTemplate.Section>

        <EuiPageTemplate.Section>
            <EuiFlexGroup>
                <EuiFlexItem grow={2}>
                    <EuiText>Docker Image Name from Docker Hub:</EuiText>
                </EuiFlexItem>
                <EuiFlexItem grow={8}>
                    <EuiFieldText
                        placeholder="Exact reference from Docker Hub"
                    >
                    </EuiFieldText>
                </EuiFlexItem>
            </EuiFlexGroup>
        </EuiPageTemplate.Section>

        <EuiPageTemplate.Section>
            <EuiButton>Create</EuiButton>
        </EuiPageTemplate.Section>
    </>)
}

export default CreateImagePage;