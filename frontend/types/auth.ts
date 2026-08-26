export interface User {
    id:number;
    username : string;
    email : string;
    avatar_url ? : string ;
    created_at?: string; 
    updated_at : string  ;
}

export interface LoginRequest  { 
    email :string ; 
    password : string;
}

export interface SignupRequest { 
    username : string;
    email : string; 
    password : string ;
    avatar_url? :string ;
}

export interface SignupRequest  { 
    username : string ; 
    email : string ; 
    password : string ; 
    avatar_url?: string ; 
}

export interface LoginResponse  { 
    token : string ;
}

export interface SignupResponse { 

    message : string ;
}
export interface ValidateResponse {
  message: string;
  user: {
    ID?: number;
    id?: number;
    Username?: string;
    username?: string;
    Email?: string;
    email?: string;
    AvatarURL?: string;
    avatar_url?: string;
    CreatedAt?: string;
    created_at?: string;
  };
}
