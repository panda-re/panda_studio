import {copyToClipboard, EuiButton, EuiButtonEmpty, EuiFieldText, EuiModal, EuiModalBody, EuiModalFooter, EuiModalHeader, EuiModalHeaderTitle, EuiOverlayMask, EuiPageTemplate, EuiSpacer, EuiText, EuiToolTip} from '@elastic/eui';
import {ReactElement, useState} from "react";
import {EuiFlexGroup, EuiFlexItem} from '@elastic/eui';
import {useLocation, useNavigate} from "react-router";
import prettyBytes from 'pretty-bytes';
import { ImageFile, PandaConfig } from '../api';

function CreateImageDetailsPage() {
  const location = useLocation()
  const navigate = useNavigate()

  const [isTextCopied, setTextCopied] = useState(false);

  const [isPopoverOpen, setIsPopoverOpen] = useState(false);

  // Modal Constants
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [modalName, setModalName] = useState("");
  const [modalDesc, setModalDesc] = useState("");
  const [modalArch, setModalArch] = useState("");
  const [modalOs, setModalOs] = useState("");
  const [modalPrompt, setModalPrompt] = useState("");
  const [modalCdrom, setModalCdrom] = useState("");
  const [modalSnapshot, setModalSnapshot] = useState("");
  const [modalMemory, setModalMemory] = useState("");
  const [modalExtraArgs, setModalExtraArgs] = useState("");

  const closeModal = () => {
    setModalName("");
    setModalDesc("");
    setModalArch("");
    setModalOs("");
    setModalPrompt("");
    setModalCdrom("");
    setModalSnapshot("");
    setModalMemory("");
    setModalExtraArgs("");
    setIsModalVisible(false)
  };
  const showModal = () => {
    setModalName(location.state.item.name);
    setModalDesc(location.state.item.description);
    setModalArch(location.state.item.config.arch);
    setModalOs(location.state.item.config.os);
    setModalPrompt(location.state.item.config.prompt);
    setModalCdrom(location.state.item.config.cdrom);
    setModalSnapshot(location.state.item.config.snapshot);
    setModalMemory(location.state.item.config.memory);
    setModalExtraArgs(location.state.item.config.extraargs);
    setIsModalVisible(true);
  }

  function createUpdateImageAndReturn(){
    if(modalName=="" || modalArch=="" || modalOs=="" || modalPrompt=="" || modalMemory==""){
      alert("Please fill out all required fields")
      return;
    }
    const conf: PandaConfig = {
      qcowfilename: location.state.item.config.qcowfilename,
      arch: modalArch,
      os: modalOs,
      prompt: modalPrompt,
      cdrom: modalCdrom,
      snapshot: modalSnapshot,
      memory: modalMemory,
      extraargs: modalExtraArgs,   
    }
    var img = {
      id: location.state.item.id,
      name: modalName,
      description: modalDesc,
      config: conf,
    }
    navigate('/images', {state: {image: img}});
  }

  function CreateConfigDisplay(){ 
    return <>
        <EuiText textAlign={"center"}>
          <strong>Qcow File Name:</strong>&nbsp;&nbsp;&nbsp;&nbsp;{location.state.item.config.qcowfilename}
        </EuiText>
        <EuiSpacer size='s'></EuiSpacer>
        <EuiText textAlign={"center"}>
          <strong>Architecture:</strong>&nbsp;&nbsp;&nbsp;&nbsp;{location.state.item.config.arch}
        </EuiText>
        <EuiSpacer size='s'></EuiSpacer>
        <EuiText textAlign={"center"}>
          <strong>OS:</strong>&nbsp;&nbsp;&nbsp;&nbsp;{location.state.item.config.os}
        </EuiText>
        <EuiSpacer size='s'></EuiSpacer>
        <EuiText textAlign={"center"}>
          <strong>Expect prompt:</strong>&nbsp;&nbsp;&nbsp;&nbsp;{location.state.item.config.prompt}
        </EuiText>
        <EuiSpacer size='s'></EuiSpacer>
        <EuiText textAlign={"center"}>
          <strong>Cdrom:</strong>&nbsp;&nbsp;&nbsp;&nbsp;{location.state.item.config.cdrom}
        </EuiText>
        <EuiSpacer size='s'></EuiSpacer>
        <EuiText textAlign={"center"}>
          <strong>Snapshot:</strong>&nbsp;&nbsp;&nbsp;&nbsp;{location.state.item.config.snapshot}
        </EuiText>
        <EuiSpacer size='s'></EuiSpacer>
        <EuiText textAlign={"center"}>
          <strong>Memory:</strong>&nbsp;&nbsp;&nbsp;&nbsp;{location.state.item.config.memory}
        </EuiText>
        <EuiSpacer size='s'></EuiSpacer>
        <EuiText textAlign={"center"}>
          <strong>Extra Args:</strong>&nbsp;&nbsp;&nbsp;&nbsp;{location.state.item.config.extraargs}
        </EuiText>
    </>
  }

  function CreateModal(){
    return <EuiOverlayMask>
              <EuiModal onClose={closeModal}>
                <EuiModalHeader>
                  <EuiModalHeaderTitle>Update Image Details</EuiModalHeaderTitle>
                </EuiModalHeader>
                <EuiModalBody>
                    <EuiFlexGroup>
                      <EuiText grow size='s'>Name:</EuiText>
                      <EuiFieldText
                        placeholder="Enter Name"  
                        name="imageName"
                        isInvalid={modalName == ""}
                        value={modalName}
                        onChange={(e) => {
                          setModalName(e.target.value);
                        }}/>
                      </EuiFlexGroup>
                      <EuiFlexGroup>
                      <EuiText size='s'>Desc:</EuiText>
                      <EuiFieldText 
                      placeholder="Enter New Description"  
                      name="imageDesc" 
                      value={modalDesc}
                      onChange={(e) => {
                        setModalDesc(e.target.value);
                      }}/>
                      </EuiFlexGroup>
                      <EuiFlexGroup>
                      <EuiText size='s'>Arch:</EuiText>
                      <EuiFieldText 
                      placeholder="Enter image Architecture (required)"
                      value={modalArch}
                      isInvalid={modalArch == ""}
                      name="pandaConfigArch" 
                      onChange={(e) => {
                        setModalArch(e.target.value);
                      }}/>
                      </EuiFlexGroup>
                      <EuiFlexGroup>
                      <EuiText size='s'>Os:</EuiText>
                      <EuiFieldText 
                      placeholder="Enter image OS (required)"
                      value={modalOs}
                      isInvalid={modalOs == ""}
                      name="pandaConfigOs" 
                      onChange={(e) => {
                        setModalOs(e.target.value);
                      }}/>
                      </EuiFlexGroup>
                      <EuiFlexGroup>
                      <EuiText size='s'>Prompt:</EuiText>
                      <EuiFieldText 
                      placeholder="Enter prompt (required)"
                      value={modalPrompt}
                      isInvalid={modalPrompt == ""}
                      name="pandaConfigPrompt" 
                      onChange={(e) => {
                        setModalPrompt(e.target.value);
                      }}/>
                      </EuiFlexGroup>
                      <EuiFlexGroup>
                      <EuiText size='s'>Cdrom:</EuiText>
                      <EuiFieldText 
                      placeholder="Enter Cdrom"
                      value={modalCdrom}
                      name="pandaConfigCdrom" 
                      onChange={(e) => {
                        setModalCdrom(e.target.value);
                      }}/>
                      </EuiFlexGroup>
                      <EuiFlexGroup>
                      <EuiText size='s'>Snapshot:</EuiText>
                      <EuiFieldText 
                      placeholder="Enter Snapshot"
                      value={modalSnapshot}
                      name="pandaConfigSnapshot" 
                      onChange={(e) => {
                        setModalSnapshot(e.target.value);
                      }}/>
                      </EuiFlexGroup>
                      <EuiFlexGroup>
                      <EuiText size='s'>Memory:</EuiText>
                      <EuiFieldText 
                      placeholder="Enter memory amount (required)"
                      value={modalMemory}
                      isInvalid={modalMemory == ""}
                      name="pandaConfigMemory" 
                      onChange={(e) => {
                        setModalMemory(e.target.value);
                      }}/>
                      </EuiFlexGroup>
                      <EuiFlexGroup>
                      <EuiText size='s'>Extra Args:</EuiText>
                      <EuiFieldText 
                      placeholder="Enter Extra args"
                      value={modalExtraArgs}
                      name="pandaConfigExtraArgs"
                      onChange={(e) => {
                        setModalExtraArgs(e.target.value);
                      }}/>
                      </EuiFlexGroup>
                </EuiModalBody>
                <EuiModalFooter>
                  <EuiButton onClick={closeModal} fill>Close</EuiButton>
                  <EuiButton 
                    onClick={createUpdateImageAndReturn}
                    fill>
                      Submit</EuiButton>
                </EuiModalFooter>
              </EuiModal>
            </EuiOverlayMask>
  }

  const buttonStyle = {
    marginRight: "25px",
    marginTop: "25px"
  }

  var size = 0;
  for(var f of location.state.item.files){
    size+= (f.size != null) ? +f.size: 0;
  }

  function CreateImageFileRows(files: ImageFile[]){
    var items: ReactElement[] = [];
    for(var file of files){
      items.push(<EuiFlexGroup>
              <EuiFlexItem>
                <EuiText textAlign={"center"}>
                  <strong>ID:</strong>
                </EuiText>
                <EuiText textAlign={"center"}>
                  {file.id}
                </EuiText>
              </EuiFlexItem>
              <EuiFlexItem>
                <EuiText textAlign={"center"}>
                  <strong>Name:</strong>
                </EuiText>
                <EuiText textAlign={"center"}>
                  {file.file_name}
                </EuiText>
              </EuiFlexItem>
              <EuiFlexItem>
                <EuiText textAlign={"center"}>
                  <strong>Type:</strong>
                </EuiText>
                <EuiText textAlign={"center"}>
                  {file.file_type}
                </EuiText>
              </EuiFlexItem>
              <EuiFlexItem>
                <EuiText textAlign={"center"}>
                  <strong>File Size:</strong>
                </EuiText>
                <EuiText textAlign={"center"}>
                  {(file.size != null) ? prettyBytes(file.size, { maximumFractionDigits: 2 }) : "0"}
                </EuiText>
              </EuiFlexItem>
              <EuiFlexItem >
                <EuiText textAlign='center'>
                  <strong>Hash:</strong>
                </EuiText>
                <EuiText textAlign='center'>
                  <EuiToolTip
                    content={isTextCopied ? 'Hash copied to clipboard' : 'Copy hash'}>
                    <EuiButtonEmpty 
                      color='text'
                      onBlur={() => {
                        setTextCopied(false);
                      }}             
                      onClick={() => {
                        copyToClipboard(file.sha256 ?? "");
                        setTextCopied(true);
                      }}>
                        {file.sha256?.substring(0, 5)}...
                    </EuiButtonEmpty>
                  </EuiToolTip>
                </EuiText>
              </EuiFlexItem>
            </EuiFlexGroup>)}
    return items;
  }

  return(<>
    <EuiPageTemplate.Header pageTitle="Image Details" />
    <EuiFlexGroup>
      <EuiFlexItem grow={6}>
        <EuiPageTemplate.Section>
          <EuiFlexGroup>
            <EuiFlexItem>
              <EuiText textAlign={"center"}>
                <strong>ID:</strong>&nbsp;&nbsp;&nbsp;&nbsp;{location.state.item.id}
              </EuiText>
            </EuiFlexItem>
            <EuiFlexItem>
              <EuiText textAlign={"center"}>
                <strong>Name:</strong>&nbsp;&nbsp;&nbsp;&nbsp;{location.state.item.name}
              </EuiText>
            </EuiFlexItem>
          </EuiFlexGroup>
          <EuiSpacer></EuiSpacer>
          <EuiFlexGroup>
            <EuiFlexItem>
              <EuiText textAlign={"center"}>
                <strong>Description:</strong>&nbsp;&nbsp;&nbsp;&nbsp;{location.state.item.description}
              </EuiText>
            </EuiFlexItem>
            <EuiFlexItem>
              <EuiText textAlign={"center"}>
                <strong>Size:</strong>&nbsp;&nbsp;&nbsp;&nbsp;{prettyBytes(size, { maximumFractionDigits: 2 })}
              </EuiText>
            </EuiFlexItem>
          </EuiFlexGroup>
          <EuiSpacer size='xxl'></EuiSpacer>
          <EuiText textAlign='center'><strong><u>Image Config</u></strong></EuiText>
          {CreateConfigDisplay()}
          <EuiSpacer size='xxl'></EuiSpacer>
          <EuiText textAlign="center"><strong><u>Image Files</u></strong></EuiText>
          {CreateImageFileRows(location.state.item.files)}
        </EuiPageTemplate.Section>
      </EuiFlexItem>

      <EuiFlexItem>
        <EuiFlexGroup direction={"column"}>
          <EuiFlexItem grow={false}>
            <EuiButton 
            style={buttonStyle}
            onClick={() => {
              navigate('/createImage', {state:{item:location.state.item}})
            }}
            >Derive New Image</EuiButton>
          </EuiFlexItem>
          <EuiFlexItem grow={false}>
          <EuiButton 
              style={buttonStyle}
              onClick= {() => {
                navigate('/images', {state: {imageId: location.state.item.id}})
              }}
            >Delete Image</EuiButton>
          </EuiFlexItem>
          <EuiFlexItem grow={false}>
            <EuiButton 
            style={buttonStyle}
            onClick={showModal}
            >Update Image Info</EuiButton>
          </EuiFlexItem>
        </EuiFlexGroup>
      </EuiFlexItem>
    </EuiFlexGroup>
    {(isModalVisible) ? (CreateModal()) : null}
  </>)
}


export default CreateImageDetailsPage;