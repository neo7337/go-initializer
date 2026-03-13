export type AddonCategory = 'cache' | 'database' | 'other';

export interface AddonState {
  cache: string[];
  database: string[];
  other: string[];
  [key: string]: string[];
}

export interface FormErrors {
  moduleName?: string;
  name?: string;
  description?: string;
  goVersion?: string;
  projectType?: string;
  framework?: string;
}

export interface FormTouched {
  moduleName: boolean;
  name: boolean;
  description: boolean;
  goVersion: boolean;
  projectType: boolean;
  framework: boolean;
}
