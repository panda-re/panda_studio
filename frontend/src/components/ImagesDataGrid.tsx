import { EuiBasicTable, EuiBasicTableColumn, EuiButton, EuiButtonIcon, EuiFieldText, EuiFilePicker, EuiFlexGroup, EuiFlexItem, EuiModal, EuiModalBody, EuiModalFooter, EuiModalHeader, EuiModalHeaderTitle, EuiOverlayMask, EuiSearchBar, EuiSearchBarOnChangeArgs, EuiSelect, EuiSpacer, EuiText, RIGHT_ALIGNMENT, useGeneratedHtmlId } from '@elastic/eui';
import { useQueryClient } from '@tanstack/react-query';
import prettyBytes from 'pretty-bytes';
import React, { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { CreateImageFileRequest, CreateImageRequest, Image, ImageFile, ImageFileType, PandaConfig, useCreateImage, useCreateImageFile, useDeleteImageById, useFindAllImages, useUpdateImage } from '../api';

function ImagesDataGrid() {
  const navigate = useNavigate();
  const location = useLocation();
  const {isLoading, isError, data} = useFindAllImages();
  const queryClient = useQueryClient();
  const deleteFunction = useDeleteImageById({
    mutation: {
      onSuccess: () => queryClient.invalidateQueries(),
      onError: (response) => alert("Error deleting Image: " + response.error?.message)}});
  const updateFn = useUpdateImage({
    mutation: {
      onSuccess: () => queryClient.invalidateQueries(),
      onError: (response) => alert("Error updating image: " + response.error?.message)}});

   // File picker constants
   const createFileFn = useCreateImageFile({mutation: {onSuccess(data, variables, context) {
    setIsLoadingVisible(false);
    queryClient.invalidateQueries();
  }}})
   const filePickerId = useGeneratedHtmlId({ prefix: 'filePicker' });
   const [files, setFiles] = useState(new Array<File>);
 
   const onFileChange = (files: FileList | null) => {
     setFiles(files!.length > 0 ? Array.from(files!) : []);
   };

   // Dropdown Constants
   const archOptions = [
    { value: 'x86_64', text: 'x86_64' },
    { value: 'i386', text: 'i386' },
    { value: 'arm', text: 'arm' },
    { value: 'aarch64', text: 'aarch64' },
    { value: 'ppc', text: 'ppc' },
    { value: 'mips', text: 'mips' },
    { value: 'mipsel', text: 'mipsel' },
    { value: 'mips64', text: 'mips64' },
    ];

    const [archValue, setArchValue] = useState(archOptions[0].value);

    const basicSelectId = useGeneratedHtmlId({ prefix: 'basicSelect' });

    const onDropdownChange = (val: string) => {
      setArchValue(val);
    };
   
   ///////// Modal Constants ///////////////////
   const [isModalVisible, setIsModalVisible] = useState(false);
   const [modalName, setModalName] = useState("");
   const [modalDesc, setModalDesc] = useState("");
   const [modalOs, setModalOs] = useState("");
   const [modalPrompt, setModalPrompt] = useState("");
   const [modalCdrom, setModalCdrom] = useState("");
   const [modalSnapshot, setModalSnapshot] = useState("");
   const [modalMemory, setModalMemory] = useState("");
   const [modalExtraArgs, setModalExtraArgs] = useState("");
 
   const [isLoadingVisible, setIsLoadingVisible] = useState(false);
 
   const closeModal = () => {
     setModalName("");
     setModalDesc("");
     setModalOs("");
     setModalPrompt("");
     setModalCdrom("");
     setModalSnapshot("");
     setModalMemory("");
     setModalExtraArgs("");
     setIsModalVisible(false)
   };
   const showModal = () => {
     setIsModalVisible(true);
   }

  
  /////////// Endpoint Functions //////////////
  const deleteImage = ({itemId}: {itemId: string}) => {
    deleteFunction.mutate({imageId: itemId});
  }

  const updateImage = ({image}: {image: Image}) => {
    if(image.id == null){
      return;
    }
    const conf: PandaConfig = {
      qcowfilename: image.config?.qcowfilename,
      arch: image.config?.arch,
      os: image.config?.os,
      prompt: image.config?.prompt,
      cdrom: image.config?.cdrom,
      snapshot: image.config?.snapshot,
      memory: image.config?.memory,
      extraargs: image.config?.extraargs,      
    }
    const req: CreateImageRequest = {
      name: image.name,
      description: image.description,
      config: conf,
    };
    updateFn.mutate({data: req, imageId: image.id});
  }

  function deleteActionPress (event: React.MouseEvent, item: Image){
    deleteImage({itemId: item.id!})
    event.stopPropagation();
  }

  function createFiles(image: Image){
    var splitArray = files[0].name.split(".");
    var fileTypeString = splitArray[splitArray.length-1];
    var fileType;
    switch(fileTypeString){
      case "qcow2": {
        fileType = ImageFileType.qcow2;
        break;
      }
      case "dtb": {
        fileType = ImageFileType.dtb;
        break;
      }
      case "kernel": {
        fileType = ImageFileType.kernel;
        break;
      }
      default:{
        fileType = undefined;
        break;
      }
    }

    if(fileType == undefined){
      // closeModal();
      alert("Invalid File Type");
      deleteImage({itemId: image.id ?? ""})
      return;
    }
    const fileReq: CreateImageFileRequest = {
      file_name: files[0].name,
      file_type: fileType,
      file: files[0],
    }
    createFileFn.mutate({data: fileReq, imageId: image.id!})
    setIsLoadingVisible(true);
    closeModal();
  }

  const createFn = useCreateImage({mutation: {onSuccess(data, variables, context) {createFiles(data)},}})

  function createFile(){
    if(modalName=="" || modalOs=="" || modalPrompt=="" || modalMemory==""){
      alert("Please fill out all required fields")
      return;
    }
    const conf: PandaConfig = {
      qcowfilename: files[0].name,
      arch: archValue,
      os: modalOs,
      prompt: modalPrompt,
      cdrom: modalCdrom,
      snapshot: modalSnapshot,
      memory: modalMemory,
      extraargs: modalExtraArgs,   
    }
    const req: CreateImageRequest = {
      name: modalName,
      description: modalDesc,
      config: conf,
    };
    createFn.mutate({data: req})
  }

  //////// UI Functions ///////////
  useEffect(() => {
    if(location.state) {
      if(location.state.image){
        updateImage({image: location.state.image});
      }
      else{
        deleteImage({itemId: location.state.imageId});
      }
      window.history.replaceState({}, document.title)
    }
  }, []);

  function LoadingModal(){
    return <EuiOverlayMask>
              <EuiModal onClose={()=>{}}>
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
                      placeholder="Enter Name (required)"
                      isInvalid={modalName == ""}
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
                      <EuiSelect
                        id={basicSelectId}
                        options={archOptions}
                        value={archValue}
                        onChange={(e) => {
                          onDropdownChange(e.target.value);
                        }}
                        aria-label="Use aria labels when no actual label is in use"
                      />
                      <EuiFieldText 
                      placeholder="Enter image OS (required)"
                      isInvalid={modalOs == ""}
                      name="pandaConfigOs" 
                      onChange={(e) => {
                        setModalOs(e.target.value);
                      }}/>
                      <EuiFieldText 
                      placeholder="Enter prompt (required)"
                      isInvalid={modalPrompt == ""}
                      name="pandaConfigPrompt" 
                      onChange={(e) => {
                        setModalPrompt(e.target.value);
                      }}/>
                      <EuiFieldText 
                      placeholder="Enter Cdrom" 
                      name="pandaConfigCdrom" 
                      onChange={(e) => {
                        setModalCdrom(e.target.value);
                      }}/>
                      <EuiFieldText 
                      placeholder="Enter Snapshot"  
                      name="pandaConfigSnapshot" 
                      onChange={(e) => {
                        setModalSnapshot(e.target.value);
                      }}/>
                      <EuiFieldText 
                      placeholder="Enter memory amount (required)"
                      isInvalid={modalMemory == ""}
                      name="pandaConfigMemory" 
                      onChange={(e) => {
                        setModalMemory(e.target.value);
                      }}/>
                      <EuiFieldText 
                      placeholder="Enter Extra args"  
                      name="pandaConfigExtraArgs"
                      onChange={(e) => {
                        setModalExtraArgs(e.target.value);
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
                    }} 
                    fill>
                      Submit</EuiButton>
                </EuiModalFooter>
              </EuiModal>
            </EuiOverlayMask>
  }

  const tableColumns: EuiBasicTableColumn<Image>[] = [
    {
      field: 'id',
      name: 'Id',
    },
    {
      field: 'name',
      name: 'File Name',
    },
    {
      field: 'files',
      name: 'Size',
      render: (value: ImageFile[]) => {
        var size = 0;
        for(var f of value){
          size+= (f.size != null) ? +f.size: 0;
        }
        return prettyBytes(size, { maximumFractionDigits: 2 });
      },
    },
    {
      align: RIGHT_ALIGNMENT,
      name: 'Delete',
      render: (item: Image) => {
        return (
          <EuiButtonIcon
            onClick={(event: React.MouseEvent) => {deleteActionPress(event, item)}}
            iconType={"trash"}
          />
        );
      },
    },
  ]
  
  const getRowProps = (item: Image) => {
    const id = item.id;
    return {
      'data-test-subj': `image-row-${id}`,
      onClick: () => {
        navigate('/imageDetails', {state:{item: item}})
      },
    }
  }

  //////////// Search Bar Items //////////////
  const initialQuery = EuiSearchBar.Query.MATCH_ALL;

  const [query, setQuery] = useState(initialQuery);

  const onChange = (args: EuiSearchBarOnChangeArgs) => {
    setQuery(args.query ?? initialQuery);
  };

  const queriedItems = EuiSearchBar.Query.execute(query, data ?? []);

  ////////// UI Element //////////
  return (<>
    <EuiFlexGroup justifyContent='spaceBetween'>
      <EuiFlexItem grow={false} style={{ minWidth: 300 }}>
        <EuiSearchBar 
          box={{
            incremental: true,
          }}
          defaultQuery={initialQuery}
          onChange={onChange}/>
      </EuiFlexItem>
      <EuiFlexItem grow={false}>
          <EuiButton onClick={showModal} iconType={'plusInCircle'}>Upload Base Image</EuiButton>
        </EuiFlexItem>
    </EuiFlexGroup>
    <EuiSpacer></EuiSpacer>
    {(isError) ? (<div>Error...</div>)
    : ((isLoading) ? (<div>Loading...</div>) 
    : <EuiBasicTable
      tableCaption="Images"
      items={queriedItems ?? []}
      rowHeader="firstName"
      columns={tableColumns}
      rowProps={getRowProps}
    />)
  }
  {(isModalVisible) ? (CreateModal()) : null}
  {(isLoadingVisible) ? (LoadingModal()) : null}
  </>)
}

export default ImagesDataGrid;