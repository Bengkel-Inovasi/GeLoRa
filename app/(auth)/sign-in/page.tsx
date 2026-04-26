import SignInForm from '../../../components/auth/SignInForm';

export const metadata = { title: 'Sign In — GeLoRa' };

export default function SignInPage() {
  return (
    <div className="w-full max-w-sm bg-white rounded-2xl shadow-lg p-8">
      <div className="mb-6 text-center">
        <h1 className="text-2xl font-bold text-slate-800">GeLoRa</h1>
        <p className="text-sm text-slate-500 mt-1">Mountain Tracker Dashboard</p>
      </div>
      <SignInForm />
    </div>
  );
}
