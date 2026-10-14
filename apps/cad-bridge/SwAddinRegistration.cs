using System;
using System.Runtime.InteropServices;
using Microsoft.Win32;

namespace SheetMetalApex.Bridge
{
    public static class SwAddinRegistration
    {
        [ComRegisterFunction]
        public static void RegisterFunction(Type t)
        {
            try
            {
                // Register the add-in in the SolidWorks registry
                string keyPath = string.Format(@"SOFTWARE\SolidWorks\Addins\{0:b}", t.GUID);
                using (RegistryKey key = Registry.LocalMachine.CreateSubKey(keyPath))
                {
                    key.SetValue(null, 1);
                    key.SetValue("Title", "SheetMetal Apex B2B Bridge");
                    key.SetValue("Description", "Send Flat Pattern and Material directly to B2B Pricing Engine.");
                }
            }
            catch (Exception ex)
            {
                Console.WriteLine("Error registering Add-in: " + ex.Message);
            }
        }

        [ComUnregisterFunction]
        public static void UnregisterFunction(Type t)
        {
            try
            {
                string keyPath = string.Format(@"SOFTWARE\SolidWorks\Addins\{0:b}", t.GUID);
                Registry.LocalMachine.DeleteSubKey(keyPath);
            }
            catch (Exception ex)
            {
                Console.WriteLine("Error unregistering Add-in: " + ex.Message);
            }
        }
    }
}
