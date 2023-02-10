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
import { createImage, createImageFile, CreateImageFileRequest, CreateImageRequest, ImageFileType, PandaConfig } from '../api';
import { AxiosRequestConfig } from 'axios';

function ImagesPage() {
  // File picker constants
  const filePickerId = useGeneratedHtmlId({ prefix: 'filePicker' });
  const [files, setFiles] = useState(new Array<File>);

  const onFileChange = (files: FileList | null) => {
    setFiles(files!.length > 0 ? Array.from(files!) : []);
  };

  const renderFiles = () => {
    if (files.length > 0) {
      return (
        <ul>
          {files.map((file, i) => (
            <li key={i}>
              <strong>{file.name}</strong> ({file.size} bytes)
            </li>
          ))}
        </ul>
      );
    } else {
      return (
        <p>Add some files to see a demo of retrieving from the FileList</p>
      );
    }
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

  function createFile(){
    const conf: PandaConfig = {
      key: modalConfig,
    }
    const req: CreateImageRequest = {
      name: modalName,
      description: modalDesc,
      config: conf,
    };
    const reqConfig: AxiosRequestConfig = {
      baseURL: "http://localhost:8080/api"
    }
    if(modalName != ""){
      createImage(req, reqConfig).then((value)=> {
        if(value.data.id != null){
          for(var f of files){
            const fileReq: CreateImageFileRequest = {
              file_name: modalName,
              file_type: ImageFileType.qcow2,
              file: f,
            }
            createImageFile(value.data.id, fileReq, reqConfig);
          }
        }
        else{
          alert("File creation error");
        }
      });
    }
    else{
      alert("File name is invalid");
    }
  }

  function CreateModal(){
    return <EuiOverlayMask>
              <EuiModal onClose={closeModal}>
                <EuiModalHeader>
                  <EuiModalHeaderTitle>Add New Interaction</EuiModalHeaderTitle>
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