import {EuiButton, EuiPageTemplate, EuiText, EuiFlexGroup, EuiFlexItem, EuiFieldText, EuiSuperSelect, EuiHealth, EuiSelect, EuiFilePicker, EuiFieldNumber} from '@elastic/eui';
import React, {ChangeEvent, SetStateAction, useState} from 'react';

//image creation needs:
// - Name
// - Size definition
// - ability to select image from docker hub
// - ability to select source files from computer
function CreateImagePage() {
    
    //TODO:
    // - need to be able to get details from image when deriving new image
    //      - add functionality to 'Derive Image' page in Image Details page to support this.
    // - get current date
    // - pull image from docker hub

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

    //hook for getting uploaded file from file drop. Currently supports uploading one file  
    const [file, setFile] = useState<SetStateAction<File>>();
    function onFileChange(newFile: File){
        console.log(newFile);
        setFile(newFile);
    }

    return (<>
        <EuiPageTemplate.Header pageTitle="Create Image"/>

        <EuiPageTemplate.Section>
        <EuiFlexGroup>
            <EuiFlexItem grow={2}>
            <EuiText>Name: </EuiText>
            </EuiFlexItem>
            <EuiFlexItem grow={8}>
            <EuiFieldText
                placeholder="eg, image1"
            />
            </EuiFlexItem>
        </EuiFlexGroup>
        <br />
        <EuiFlexGroup>
            <EuiFlexItem grow={2}>
                <EuiText>Size: </EuiText>
            </EuiFlexItem>

            <EuiFlexItem grow={8}>
                <EuiFlexGroup>
                    <EuiFlexItem >
                        <EuiFieldNumber/>
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
                           
            </EuiFlexGroup>
        </EuiPageTemplate.Section>

        <EuiPageTemplate.Section>
            <EuiFlexGroup>
                <EuiFlexItem grow={2}>
                    <EuiText>Pull Image from Docker Hub:</EuiText>
                </EuiFlexItem>
                <EuiFlexItem grow={8}>
                    <EuiFieldText
                        placeholder="placeholder"
                    >
                    </EuiFieldText>
                </EuiFlexItem>
            </EuiFlexGroup>
            <br />
            <EuiFlexGroup>
                <EuiFlexItem grow={2}>
                    <EuiText>Upload VM Image file: </EuiText>
                </EuiFlexItem>
                <EuiFlexItem grow={8}>
                    <EuiFilePicker
                        multiple={false}
                        accept={".qcow2"}
                        onChange={
                            files => {if(files!=null) onFileChange(files[0]);}}
                    >
                    </EuiFilePicker>
                    <EuiText>{file!=null ? file.name : "No File Selected"}</EuiText>
                </EuiFlexItem>
                
            </EuiFlexGroup>
        </EuiPageTemplate.Section>

        <EuiPageTemplate.Section>
        <EuiButton>Create</EuiButton>
        </EuiPageTemplate.Section>

        
    
    </>)
}

export default CreateImagePage;