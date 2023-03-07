import {EuiButton, EuiFieldText, EuiFilePicker, EuiModal, EuiModalBody, EuiModalFooter, EuiModalHeader, EuiModalHeaderTitle, EuiOverlayMask, EuiPageTemplate, EuiSpacer, EuiText} from '@elastic/eui';
import {ReactElement, useState} from "react";
import {EuiFlexGroup, EuiFlexItem} from '@elastic/eui';
import {useLocation, useNavigate} from "react-router";
import prettyBytes from 'pretty-bytes';
import { CreateImageRequest, ImageFile, PandaConfig, useDeleteImageById, useUpdateImage } from '../api';

function CreateImageDetailsPage() {
  const location = useLocation()
  const navigate = useNavigate()
  

  // Modal Constants
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [modalName, setModalName] = useState("");
  const [modalDesc, setModalDesc] = useState("");
  const [modalConfig, setModalConfig] = useState("");

  const closeModal = () => {
    setModalName("");
    setModalDesc("");
    setModalConfig("");
    setIsModalVisible(false)
  };
  const showModal = () => {
    setModalName(location.state.item.name);
    setModalDesc(location.state.item.description);
    setModalConfig(location.state.item.config);
    setIsModalVisible(true);
  }

  function createUpdateImageAndReturn(){
    var img = {
      id: location.state.item.id,
      name: modalName,
      description: modalDesc,
      config: modalConfig,
    }
    navigate('/images', {state: {image: img}});
  }

  function CreateModal(){
    return <EuiOverlayMask>
              <EuiModal onClose={closeModal}>
                <EuiModalHeader>
                  <EuiModalHeaderTitle>Update Image Details</EuiModalHeaderTitle>
                </EuiModalHeader>
                <EuiModalBody>
                    <EuiFieldText 
                      placeholder="Enter Name"  
                      name="imageName"
                      value={modalName}
                      onChange={(e) => {
                        setModalName(e.target.value);
                      }}/>
                      <EuiFieldText 
                      placeholder="Enter New Description"  
                      name="imageDesc" 
                      value={modalDesc}
                      onChange={(e) => {
                        setModalDesc(e.target.value);
                      }}/>
                      <EuiFieldText 
                      placeholder="Enter New Config"  
                      name="imageConfig" 
                      value={modalConfig}
                      onChange={(e) => {
                        setModalConfig(e.target.value);
                      }}/>
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
            </EuiFlexGroup>)
    }
    return items;
  }

  return(<>
    <EuiPageTemplate.Header pageTitle="Image Details" />
    <EuiFlexGroup>
      <EuiFlexItem grow={6}>
        <EuiPageTemplate.Section>
          <EuiText textAlign={"center"}>
            <strong>ID:</strong>
          </EuiText>
          <EuiText textAlign={"center"}>
            {location.state.item.id}
          </EuiText>
          <EuiSpacer></EuiSpacer>
          <EuiText textAlign={"center"}>
            <strong>Name:</strong>
          </EuiText>
          <EuiText textAlign={"center"}>
            {location.state.item.name}
          </EuiText>
          <EuiSpacer></EuiSpacer>
          <EuiText textAlign={"center"}>
            <strong>Description:</strong>
          </EuiText>
          <EuiText textAlign={"center"}>
            {location.state.item.description}
          </EuiText>
          <EuiSpacer></EuiSpacer>
          <EuiText textAlign={"center"}>
            <strong>Size:</strong>
          </EuiText>
          <EuiText textAlign={"center"}>
            {prettyBytes(size, { maximumFractionDigits: 2 })}
          </EuiText>
          <EuiSpacer></EuiSpacer>
        </EuiPageTemplate.Section>
        <EuiPageTemplate.Section>
          <EuiText textAlign={"center"}><strong>Image Files</strong></EuiText>
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