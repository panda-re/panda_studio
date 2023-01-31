export interface Recording {
  id: string;
  name: string;
  date: Date;
  imageName: string;
  size: number;
};

// image ID, name, OS, Timestamp, Size, view specs
export interface Image {
  id: string;
  name: string;
  date: Date;
  operatingSystem: string;
  size: number;
};

export interface InteractionProgram {
  id: string;
  name: string;
  date: Date;
};