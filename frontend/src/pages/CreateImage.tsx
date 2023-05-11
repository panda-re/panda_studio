import {EuiButton, EuiPageTemplate, EuiText, EuiFlexGroup, EuiFlexItem, EuiFieldText, EuiSuperSelect, EuiHealth, EuiSelect, EuiFilePicker, EuiFieldNumber} from '@elastic/eui';
import React, {ChangeEvent, SetStateAction, useState} from 'react';
import { useLocation } from 'react-router';
import {DeriveImageFileRequest, useCreateDerivedImage} from '../api';

function CreateImagePage() {
    const location = useLocation()
    const [newImageName, setNewImageName] = useState("");
    const [newSize, setNewSize] = useState("");
    const [dockerImageName, setDockerImageName] = useState("");

    //console.log(location.state.item);

    const deriveImageFileFn = useCreateDerivedImage({mutation: {onSuccess(data, variables, context) {console.log("derivation complete")}}})

    function sendDeriveImageRequest(){
        const req: DeriveImageFileRequest = {
            imageid: location.state.item.id,
            oldname: location.state.item.name,
            newname: newImageName,
            dockerhubimagename: dockerImageName,
            size: ""+newSize+"G" 
        };
        const fullReq = {
            data: req,
            imageId: location.state.item.id,
        };
        deriveImageFileFn.mutate(fullReq);
        console.log(fullReq);
    }

    return (<>
        <EuiPageTemplate.Header pageTitle={"Derive Image from " + location.state.item.name}/>

        <EuiPageTemplate.Section>
            <EuiFlexGroup>
                <EuiFlexItem grow={2}>
                    <EuiText>New Image name: </EuiText>
                </EuiFlexItem>
                <EuiFlexItem grow={8}>
                    <EuiFieldText
                        placeholder="eg, image1"
                        name="newname"
                        onChange={(e) => {
                            setNewImageName(e.target.value);
                        }}
                    />
                </EuiFlexItem>
            </EuiFlexGroup>
        </EuiPageTemplate.Section>
        
        <EuiPageTemplate.Section>
            <EuiFlexGroup>
                <EuiFlexItem grow={2}>
                    <EuiText>Expand to Size: (GB)</EuiText>
                </EuiFlexItem>

                <EuiFlexItem grow={8}>
                    <EuiFlexGroup>
                        <EuiFlexItem >
                            <EuiFieldNumber
                            name="size"
                            onChange={(e) => {
                                setNewSize(e.target.value);
                            }}
                            />
                        </EuiFlexItem>
                    </EuiFlexGroup>
                </EuiFlexItem>     
            </EuiFlexGroup>
        </EuiPageTemplate.Section>

        <EuiPageTemplate.Section>
            <EuiFlexGroup>
                <EuiFlexItem grow={2}>
                    <EuiText>Docker Hub image name:</EuiText>
                </EuiFlexItem>
                <EuiFlexItem grow={8}>
                    <EuiFieldText
                        placeholder="Exact reference from Docker Hub"
                        name="dockerhubimagename"
                        onChange={(e) => {
                            setDockerImageName(e.target.value);
                        }}
                    >
                    </EuiFieldText>
                </EuiFlexItem>
            </EuiFlexGroup>
        </EuiPageTemplate.Section>

        <EuiPageTemplate.Section>
            <EuiButton
            onClick={sendDeriveImageRequest}
            >Begin Derive Job</EuiButton>
        </EuiPageTemplate.Section>
    </>)

    
}



export default CreateImagePage;