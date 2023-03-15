import {
  EuiButton,
  EuiFieldSearch,
  EuiFieldText,
  EuiFilePicker,
  EuiFlexItem,
  EuiModal,
  EuiModalBody,
  EuiModalFooter,
  EuiModalHeader,
  EuiModalHeaderTitle,
  EuiOverlayMask,
  EuiPageTemplate,
  EuiSpacer,
  EuiText,
  useGeneratedHtmlId
} from '@elastic/eui';
import ImagesDataGrid from '../components/ImagesDataGrid';
import {EuiFlexGrid} from "@elastic/eui";
import { useState } from 'react';
import { CreateImageFileRequest, CreateImageRequest, Image, ImageFileType, PandaConfig, useCreateImage, useCreateImageFile } from '../api';
import { useNavigate } from 'react-router';

function ImagesPage() {
  const navigate = useNavigate()
  // File picker constants
  const createFileFn = useCreateImageFile({mutation: {onSuccess(data, variables, context) {setIsLoadingVisible(false)}}})
  const filePickerId = useGeneratedHtmlId({ prefix: 'filePicker' });
  const [files, setFiles] = useState(new Array<Blob>);

  const onFileChange = (files: FileList | null) => {
    setFiles(files!.length > 0 ? Array.from(files!) : []);
  };
  
  // Modal Constants
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [modalName, setModalName] = useState("");
  const [modalDesc, setModalDesc] = useState("");
  const [modalConfig, setModalConfig] = useState("");

  const [isLoadingVisible, setIsLoadingVisible] = useState(false);

  const closeModal = () => {
    setModalName("");
    setModalDesc("");
    setModalConfig("");
    setIsModalVisible(false)
  };
  const showModal = () => {
    setIsModalVisible(true);
  }

  function createFiles(image: Image){
    const fileReq: CreateImageFileRequest = {
      file_name: modalName,
      file_type: ImageFileType.qcow2,
      file: files[0],
    }
    createFileFn.mutate({data: fileReq, imageId: image.id!})
    setIsLoadingVisible(true);
    closeModal();
  }

  const createFn = useCreateImage({mutation: {onSuccess(data, variables, context) {createFiles(data)},}})

  function createFile(){
    const conf: PandaConfig = {
      key: modalConfig,
    }
    const req: CreateImageRequest = {
      name: modalName,
      description: modalDesc,
      config: conf,
    };
    if(modalName != ""){
      createFn.mutate({data: req})
    }
  }

  function LoadingModal(){
    return <EuiOverlayMask>
              <EuiModal onClose={closeModal}>
                <EuiModalHeader>
                  <EuiModalHeaderTitle>Uploading Image</EuiModalHeaderTitle>
                </EuiModalHeader>
                <EuiModalBody>
                    <EuiText>
                      Loading...
                    </EuiText>
                </EuiModalBody>
              </EuiModal>
            </EuiOverlayMask>
  }

  function CreateModal(){
    return <EuiOverlayMask>
              <EuiModal onClose={closeModal}>
                <EuiModalHeader>
                  <EuiModalHeaderTitle>Upload New Image</EuiModalHeaderTitle>
                </EuiModalHeader>
                <EuiModalBody>
                    <EuiFieldText 
                      placeholder="Enter Name"  
                      name="imageName" 
                      onChange={(e) => {
                        setModalName(e.target.value);
                      }}/>
                      <EuiFieldText 
                      placeholder="Enter New Description"  
                      name="imageDesc" 
                      onChange={(e) => {
                        setModalDesc(e.target.value);
                      }}/>
                      <EuiFieldText 
                      placeholder="Enter config key"  
                      name="pandaConfig" 
                      onChange={(e) => {
                        setModalConfig(e.target.value);
                      }}/>
                      <EuiFilePicker
                        id={filePickerId}
                        initialPromptText="Select or drag and drop multiple files"
                        onChange={onFileChange}
                        aria-label="Use aria labels when no actual label is in use"
                      />
                </EuiModalBody>
                <EuiModalFooter>
                  <EuiButton onClick={closeModal} fill>Close</EuiButton>
                  <EuiButton 
                    onClick={() => {
                      createFile();
                      // closeModal();
                    }} 
                    fill>
                      Submit</EuiButton>
                </EuiModalFooter>
              </EuiModal>
            </EuiOverlayMask>
  }

  return (<>
    <EuiPageTemplate.Header pageTitle='Image Dashboard' rightSideItems={[]} />

    <EuiPageTemplate.Section>
      <EuiFlexGrid columns={4}>
        {/* <EuiFlexItem>
          <EuiFieldSearch
            placeholder="Enter Image ID"
          />
        </EuiFlexItem>
        <EuiFlexItem>
          <EuiFieldSearch
            placeholder="Enter Image Name"
          />
        </EuiFlexItem>
        <EuiFlexItem>
          <EuiFieldSearch
            placeholder="Enter Date"
          />
        </EuiFlexItem> */}
        <EuiFlexItem>
          <EuiButton onClick={showModal} iconType={'plusInCircle'}>Upload Base Image</EuiButton>
        </EuiFlexItem>
      </EuiFlexGrid>
      <EuiSpacer size="xl" />
      <ImagesDataGrid />
      {(isModalVisible) ? (CreateModal()) : null}
      {(isLoadingVisible) ? (LoadingModal()) : null}
    </EuiPageTemplate.Section>

  </>)
}

export default ImagesPage;