import {EuiButton, EuiPageTemplate, EuiSpacer, EuiText} from '@elastic/eui';
import {ReactElement} from "react";
import {EuiFlexGroup, EuiFlexItem} from '@elastic/eui';
import {useLocation, useNavigate} from "react-router";
import prettyBytes from 'pretty-bytes';
import { ImageFile } from '../api';

function CreateImageDetailsPage() {
  const location = useLocation()
  const navigate = useNavigate()

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
            <EuiButton style={buttonStyle}>Delete Image</EuiButton>
          </EuiFlexItem>
        </EuiFlexGroup>
      </EuiFlexItem>
    </EuiFlexGroup>
  </>)
}


export default CreateImageDetailsPage;