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
  useGeneratedHtmlId
} from '@elastic/eui';
import ImagesDataGrid from '../components/ImagesDataGrid';
import {EuiFlexGrid} from "@elastic/eui";
import { useState } from 'react';
import { CreateImageFileRequest, CreateImageRequest, Image, ImageFileType, PandaConfig, useCreateImage, useCreateImageFile } from '../api';

function ImagesPage() {
  const createFileFn = useCreateImageFile()
  // File picker constants
  const filePickerId = useGeneratedHtmlId({ prefix: 'filePicker' });
  const [files, setFiles] = useState(new Array<File>);

  const onFileChange = (files: FileList | null) => {
    setFiles(files!.length > 0 ? Array.from(files!) : []);
  };
  
  // Modal Constants
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [modalName, setModalName] = useState("");
  const [modalDesc, setModalDesc] = useState("");
  const [modalConfig, setModalConfig] = useState("");
  const [modalType, setModalType] = useState("");

  const closeModal = () => {
    setModalName("");
    setModalDesc("");
    setModalConfig("");
    setModalType("");
    setIsModalVisible(false)
  };
  const showModal = () => {
    setIsModalVisible(true);
  }

  function createImageFiles(image: Image){
    for(var f of files){
      const fileReq: CreateImageFileRequest = {
        file_name: modalName,
        file_type: ImageFileType.qcow2,
        file: f,
      }
      createFileFn.mutate({data: fileReq, imageId: image.id ?? ""})
    }
  }

  const createFn = useCreateImage({mutation: {onSuccess(data, variables, context) {
    createImageFiles(data);
  },}})

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
                      placeholder="Enter File Type"  
                      name="imageFileType" 
                      onChange={(e) => {
                        setModalType(e.target.value);
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
                      closeModal();
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
        <EuiFlexItem>
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
        </EuiFlexItem>
        <EuiFlexItem>
          <EuiButton onClick={showModal} iconType={'plusInCircle'}>Upload Base Image</EuiButton>
        </EuiFlexItem>
      </EuiFlexGrid>
      <EuiSpacer size="xl" />
      <ImagesDataGrid />
      {(isModalVisible) ? (CreateModal()) : null}
    </EuiPageTemplate.Section>

  </>)
}

export default ImagesPage;